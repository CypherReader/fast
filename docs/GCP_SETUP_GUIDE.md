# Step-by-Step Guide: Google Cloud Platform Project Setup

## Overview

This guide will walk you through creating a GCP project, enabling required services, and setting up the infrastructure for FastingHero UAT deployment.

**Estimated Time:** 30-45 minutes  
**Cost:** Free tier available, ~$20-35/month after deployment

---

## Prerequisites

- [ ] Google account (Gmail)
- [ ] Credit card (required for GCP, even with free tier)
- [ ] Web browser

---

## Step 1: Create GCP Account & Project

### 1.1 Sign Up for Google Cloud

1. Go to <https://console.cloud.google.com>
2. Click **"Get started for free"** or **"Go to Console"**
3. Sign in with your Google account
4. Accept the Terms of Service

**First-time users get:**

- $300 free credits (valid for 90 days)
- Always-free tier for many services

### 1.2 Enter Billing Information

1. Click **"Activate"** or **"Enable Billing"**
2. Select your country
3. Enter credit card details
   - ⚠️ You won't be charged during free trial
   - After trial: You must manually upgrade to paid
4. Click **"Start my free trial"**

### 1.3 Create New Project

1. In the top bar, click the project dropdown (next to "Google Cloud")
2. Click **"New Project"**
3. Enter project details:

   ```
   Project name: FastingHero UAT
   Organization: (leave as "No organization" if personal account)
   Location: (leave default)
   ```

4. Click **"CREATE"**
5. Wait 10-30 seconds for project creation
6. **Copy your Project ID** (e.g., `fastinghero-uat-123456`)
   - You'll need this later!

---

## Step 2: Enable Required APIs

### 2.1 Navigate to APIs & Services

1. Open the hamburger menu (☰) in top-left
2. Click **"APIs & Services"** → **"Library"**

### 2.2 Enable Each API (One by One)

Search for and enable these APIs:

#### ✅ Cloud Run API

1. Search: "Cloud Run API"
2. Click the result
3. Click **"ENABLE"**
4. Wait for confirmation

#### ✅ Cloud SQL Admin API

1. Search: "Cloud SQL Admin API"
2. Click **"ENABLE"**

#### ✅ Secret Manager API

1. Search: "Secret Manager API"
2. Click **"ENABLE"**

#### ✅ Artifact Registry API (replaces Container Registry)

1. Search: "Artifact Registry API"
2. Click **"ENABLE"**

#### ✅ Cloud Build API

1. Search: "Cloud Build API"
2. Click **"ENABLE"**

#### ✅ Compute Engine API (for Cloud SQL)

1. Search: "Compute Engine API"
2. Click **"ENABLE"**

**Verification:**

- Go to **"APIs & Services"** → **"Dashboard"**
- You should see 6 enabled APIs

---

## Step 3: Set Up Cloud SQL Database

### 3.1 Navigate to Cloud SQL

1. Hamburger menu (☰) → **"SQL"**
2. Click **"Create Instance"**

### 3.2 Choose PostgreSQL

1. Click **"Choose PostgreSQL"**
2. Click **"Enable Compute Engine API"** if prompted

### 3.3 Configure Instance

**Instance Info:**

```
Instance ID: fastinghero-uat-db
Password: [Generate strong password - save it!] - Curacao2020!!
```

**Choose Configuration:**

- Click **"SHOW CONFIGURATION OPTIONS"**

**Database version:**

```
PostgreSQL 15
```

**Choose region and zonal availability:**

```
Region: us-central1 (Iowa)
  - Or choose closest to your users
Zonal availability: Single zone
```

**Customize your instance:**

**Machine Configuration:**

```
Machine type: Shared core
  → Cores: 1 vCPU
  → Memory: 0.614 GB
```

**Storage:**

```
Storage type: SSD
Storage capacity: 10 GB
Enable automatic storage increases: ✓ (checked)
```

**Connections:**

```
Instance IP assignment:
  ☑ Public IP
  ☑ Private IP (recommended for production, optional for UAT)
```

**Data Protection:**

```
Automate backups: ✓ (checked)
Backup window: (leave default or choose off-peak)
Point-in-time recovery: ✓ (checked)
```

**Maintenance:**

```
Leave defaults (allows automatic updates)
```

### 3.4 Create Instance

1. Click **"CREATE INSTANCE"**
2. **This takes 5-10 minutes** - get coffee! ☕
3. Save the following information:

```
Instance Connection Name: [PROJECT_ID]:us-central1:fastinghero-uat-db
Public IP Address: [will appear after creation]
Root Password: [the one you generated]
```

---

## Step 4: Create Database

### 4.1 Connect to Cloud SQL

1. Once instance is created, click its name
2. Go to **"Databases"** tab
3. Click **"CREATE DATABASE"**

### 4.2 Database Configuration

```
Database name: fastinghero
Character set: UTF8 (default)
Collation: en_US.UTF8 (default)
```

4. Click **"CREATE"**

---

## Step 5: Set Up Secret Manager

### 5.1 Navigate to Secret Manager

1. Hamburger menu (☰) → **"Security"** → **"Secret Manager"**
2. Click **"CREATE SECRET"**

### 5.2 Create Secrets (Repeat for Each)

#### Secret 1: JWT_SECRET

```
Name: JWT_SECRET
Secret value: [Generate random 48-character string]
```

**Generate strong secret in PowerShell:**

```powershell
-join ((65..90) + (97..122) + (48..57) | Get-Random -Count 48 | ForEach-Object {[char]$_})
```

Click **"CREATE SECRET"**

#### Secret 2: DATABASE_URL

```
Name: DATABASE_URL
Secret value: postgres://postgres:[YOUR_DB_PASSWORD]@[PUBLIC_IP]:5432/fastinghero?sslmode=require
```

Replace:

- `[YOUR_DB_PASSWORD]` with the password from Step 3.3
- `[PUBLIC_IP]` with your Cloud SQL public IP

Click **"CREATE SECRET"**

#### Secret 3: GOOGLE_CLIENT_ID (We'll update this later)

```
Name: GOOGLE_CLIENT_ID
Secret value: placeholder-will-update-after-oauth-setup
```

#### Secret 4: GOOGLE_CLIENT_SECRET (We'll update this later)

```
Name: GOOGLE_CLIENT_SECRET
Secret value: placeholder-will-update-after-oauth-setup
```

#### Optional Secret 5: DEEPSEEK_API_KEY

```
Name: DEEPSEEK_API_KEY
Secret value: [Your DeepSeek API key]
```

---

## Step 6: Set Up Service Account

### 6.1 Create Service Account

1. Hamburger menu (☰) → **"IAM & Admin"** → **"Service Accounts"**
2. Click **"CREATE SERVICE ACCOUNT"**

### 6.2 Service Account Details

```
Service account name: fastinghero-uat-runner
Service account ID: fastinghero-uat-runner (auto-filled)
Description: Service account for Cloud Run deployment
```

Click **"CREATE AND CONTINUE"**

### 6.3 Grant Permissions

Add these roles:

1. Click **"Select a role"**
2. Add: **"Cloud Run Admin"**
3. Click **"Add Another Role"**
4. Add: **"Cloud SQL Client"**
5. Click **"Add Another Role"**
6. Add: **"Secret Manager Secret Accessor"**

Click **"CONTINUE"** → **"DONE"**

### 6.4 Create Service Account Key (for CI/CD)

1. Click on the service account you just created
2. Go to **"KEYS"** tab
3. Click **"ADD KEY"** → **"Create new key"**
4. Choose **"JSON"**
5. Click **"CREATE"**
6. **Save this file securely** - you'll need it for GitHub Actions

---

## Step 7: Configure Google OAuth

### 7.1 Go to OAuth Consent Screen

1. Hamburger menu (☰) → **"APIs & Services"** → **"OAuth consent screen"**
2. Choose **"External"** (allows any Google account)
3. Click **"CREATE"**

### 7.2 App Information

```
App name: FastingHero UAT
User support email: [Your email]
App logo: (optional - skip for now)
Application home page: (leave empty for now)
Application privacy policy: (leave empty for now)
Application terms of service: (leave empty for now)
Authorized domains: (leave empty for now - add after deployment)
Developer contact email: [Your email]
```

Click **"SAVE AND CONTINUE"**

### 7.3 Scopes

1. Click **"ADD OR REMOVE SCOPES"**
2. Select these scopes:
   - ✓ `.../auth/userinfo.email`
   - ✓ `.../auth/userinfo.profile`
   - ✓ `openid`
3. Click **"UPDATE"**
4. Click **"SAVE AND CONTINUE"**

### 7.4 Test Users (while app is in testing)

1. Click **"ADD USERS"**
2. Add your email and any test user emails
3. Click **"SAVE AND CONTINUE"**

Click **"BACK TO DASHBOARD"**

### 7.5 Create OAuth 2.0 Credentials

1. Go to **"Credentials"** tab (left sidebar)
2. Click **"CREATE CREDENTIALS"** → **"OAuth client ID"**
3. Application type: **"Web application"**
4. Name: `FastingHero UAT Web Client`

**Authorized JavaScript origins:**

```
http://localhost:5173
(Will add Cloud Run URL later)
```

**Authorized redirect URIs:**

```
http://localhost:5173/auth/callback
http://localhost:8080/api/v1/auth/google/callback
(Will add Cloud Run URLs later)
```

5. Click **"CREATE"**
6. **Copy and save:**
   - Client ID - 537397575496-k0vhv9gq1dr2f0mu0ln28l5fnlobuhbh.apps.googleusercontent.com
   - Client Secret

### 7.6 Update Secrets

1. Go back to **Secret Manager**
2. Update `GOOGLE_CLIENT_ID`:
   - Click the secret → **"NEW VERSION"**
   - Paste your Client ID → **"ADD NEW VERSION"**
3. Update `GOOGLE_CLIENT_SECRET`:
   - Click the secret → **"NEW VERSION"**
   - Paste your Client Secret → **"ADD NEW VERSION"**

---

## Step 8: Install Google Cloud CLI (on your computer)

### 8.1 Download and Install

**Windows:**

1. Download: <https://cloud.google.com/sdk/docs/install>
2. Run the installer
3. Choose default options
4. Check "Run `gcloud init`" at the end

**Alternative - PowerShell:**

```powershell
(New-Object Net.WebClient).DownloadFile("https://dl.google.com/dl/cloudsdk/channels/rapid/GoogleCloudSDKInstaller.exe", "$env:Temp\GoogleCloudSDKInstaller.exe")
& $env:Temp\GoogleCloudSDKInstaller.exe
```

### 8.2 Initialize gcloud

```bash
gcloud init
```

Follow prompts:

1. Choose: "Log in with a new account"
2. Browser will open - sign in
3. Choose your project: **FastingHero UAT**
4. Choose region: **us-central1-a** (or your Cloud SQL region)

### 8.3 Set Default Project

```bash
gcloud config set project [YOUR_PROJECT_ID]
```

---

## Step 9: Verify Setup

### 9.1 Check Enabled APIs

```bash
gcloud services list --enabled
```

You should see all 6 APIs enabled.

### 9.2 Check Cloud SQL Instance

```bash
gcloud sql instances list
```

You should see your `fastinghero-uat-db` instance.

### 9.3 Check Secrets

```bash
gcloud secrets list
```

You should see 4-5 secrets.

---

## ✅ Checklist: What You Should Have Now

- [ ] GCP Project created (Project ID saved)
- [ ] Billing enabled
- [ ] 6 APIs enabled (Cloud Run, Cloud SQL, Secret Manager, Artifact Registry, Cloud Build, Compute Engine)
- [ ] Cloud SQL PostgreSQL instance running
- [ ] Database `fastinghero` created
- [ ] 5 secrets in Secret Manager (JWT_SECRET, DATABASE_URL, GOOGLE_CLIENT_ID, GOOGLE_CLIENT_SECRET, DEEPSEEK_API_KEY)
- [ ] Service account `fastinghero-uat-runner` created with 3 roles
- [ ] Service account JSON key downloaded
- [ ] OAuth consent screen configured
- [ ] OAuth 2.0 credentials created (Client ID and Secret)
- [ ] gcloud CLI installed and authenticated

---

## 📝 Save This Information

Create a secure note with:

```
GCP Project ID: ___________________
Cloud SQL Instance Name: ___________________
Cloud SQL Public IP: ___________________
Cloud SQL Password: ___________________
Google OAuth Client ID: ___________________
Google OAuth Client Secret: ___________________
Service Account Key: [JSON file location]
```

---

## Next Steps

Now that your GCP project is set up, you're ready to:

1. ✅ **Build Docker image** for your application
2. ✅ **Deploy to Cloud Run**
3. ✅ **Implement Google OAuth** in backend
4. ✅ **Set up CI/CD pipeline**

Would you like me to proceed with the next phase?

---

## Troubleshooting

### "Billing account required"

- Go to **Billing** in the menu
- Link a billing account
- You won't be charged during free trial

### "Quota exceeded"

- Free tier has limits
- Check **IAM & Admin** → **Quotas**
- Request quota increase if needed

### "API not enabled"

- Go to **APIs & Services** → **Library**
- Search for the API
- Click **ENABLE**

### Cloud SQL creation fails

- Ensure Compute Engine API is enabled
- Check billing is active
- Try different region

---

## Cost Monitoring

Set up budget alerts to avoid surprises:

1. Go to **Billing** → **Budgets & alerts**
2. Click **CREATE BUDGET**
3. Set budget: $50/month
4. Set alert threshold: 50%, 90%, 100%
5. Add your email for notifications

---

**Setup Complete!** 🎉

You now have a production-ready GCP environment for FastingHero UAT deployment.
