# Deployment Guide: Google Cloud Run

## Quick Deployment Steps

### 1. Prerequisites Completed

Before deploying, ensure you have completed all steps in [`docs/GCP_SETUP_GUIDE.md`](./GCP_SETUP_GUIDE.md):

- ✅ GCP Project created
- ✅ APIs enabled
- ✅ Cloud SQL database running
- ✅ Secrets configured in Secret Manager
- ✅ Service account created with JSON key
- ✅ Google OAuth credentials created

### 2. Update Configuration Files

#### Update `cloudbuild.yaml`

Replace the placeholder in `cloudbuild.yaml`:

```yaml
substitutions:
  _CLOUD_SQL_INSTANCE: '[YOUR_PROJECT_ID]:[REGION]:fastinghero-uat-db'
```

Example:

```yaml
substitutions:
  _CLOUD_SQL_INSTANCE: 'fastinghero-uat-123456:us-central1:fastinghero-uat-db'
```

#### Update GitHub Secrets (if using GitHub Actions)

Add these secrets to your GitHub repository (Settings → Secrets and variables → Actions):

```
GCP_PROJECT_ID: your-project-id
GCP_SA_KEY: [paste service account JSON key]
CLOUD_SQL_INSTANCE: project-id:region:instance-name
```

### 3. Manual Deployment (First Time)

```bash
# Set your project ID
export PROJECT_ID="your-project-id-here"

# Build the Docker image
docker build -t gcr.io/$PROJECT_ID/fastinghero-uat:latest .

# Configure Docker to use gcloud
gcloud auth configure-docker

# Push to Google Container Registry
docker push gcr.io/$PROJECT_ID/fastinghero-uat:latest

# Deploy to Cloud Run
gcloud run deploy fastinghero-uat \
  --image gcr.io/$PROJECT_ID/fastinghero-uat:latest \
  --region us-central1 \
  --platform managed \
  --allow-unauthenticated \
  --set-env-vars ENVIRONMENT=uat \
  --set-secrets JWT_SECRET=JWT_SECRET:latest,DATABASE_URL=DATABASE_URL:latest,GOOGLE_CLIENT_ID=GOOGLE_CLIENT_ID:latest,GOOGLE_CLIENT_SECRET=GOOGLE_CLIENT_SECRET:latest,DEEPSEEK_API_KEY=DEEPSEEK_API_KEY:latest \
  --add-cloudsql-instances YOUR_PROJECT:REGION:fastinghero-uat-db \
  --memory 512Mi \
  --cpu 1 \
  --timeout 300 \
  --max-instances 10 \
  --min-instances 0
```

### 4. Get Your Deployment URL

After deployment:

```bash
gcloud run services describe fastinghero-uat --region us-central1 --format 'value(status.url)'
```

Example output: `https://fastinghero-uat-abc123-uc.a.run.app`

### 5. Update Google OAuth Redirect URIs

1. Go to [Google Cloud Console → APIs & Services → Credentials](https://console.cloud.google.com/apis/credentials)
2. Click on your OAuth 2.0 Client ID
3. Add your Cloud Run URL to **Authorized redirect URIs**:

   ```
   https://your-cloud-run-url/api/v1/auth/google/callback
   ```

4. Add to **Authorized JavaScript origins**:

   ```
   https://your-cloud-run-url
   ```

5. Click **SAVE**

### 6. Test Your Deployment

```bash
# Health check
curl https://your-cloud-run-url/health

# Test registration (should fail with strong password requirement)
curl -X POST https://your-cloud-run-url/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"weak","name":"Test"}'

# Test Google OAuth (in browser)
# Navigate to: https://your-cloud-run-url/api/v1/auth/google
```

### 7. Run Database Migrations

```bash
# Connect to Cloud SQL using Cloud SQL Proxy
cloud-sql-proxy PROJECT_ID:REGION:fastinghero-uat-db &

# Run migrations
psql "host=127.0.0.1 port=5432 dbname=fastinghero user=postgres" < migrations/01_init.sql
psql "host=127.0.0.1 port=5432 dbname=fastinghero user=postgres" < migrations/02_add_timestamps.sql
psql "host=127.0.0.1 port=5432 dbname=fastinghero user=postgres" < migrations/03_add_oauth_fields.sql
```

---

## Automated Deployment with GitHub Actions

### Setup

1. Push code to GitHub
2. Add secrets as described in step 2
3. Push to `main` or `uat` branch

GitHub Actions will automatically:

- Build Docker image
- Push to GCR
- Deploy to Cloud Run
- Output the service URL

### Triggering Manual Deployment

```bash
# From GitHub UI: Actions → Deploy to GCP UAT → Run workflow

# Or push a tag
git tag uat-v1.0.0
git push origin uat-v1.0.0
```

---

## Monitoring & Logs

### View Logs

```bash
# Real-time logs
gcloud run services logs tail fastinghero-uat --region us-central1

# Last 100 lines
gcloud run services logs read fastinghero-uat --region us-central1 --limit 100
```

### View in Console

1. Go to [Cloud Run](https://console.cloud.google.com/run)
2. Click on `fastinghero-uat`
3. Go to **LOGS** tab

### Metrics

1. Cloud Run console → `fastinghero-uat` → **METRICS** tab
2. View:
   - Request count
   - Request latency
   - Error rate
   - CPU/Memory usage

---

## Troubleshooting

### "Secret not found"

- Verify secrets exist in Secret Manager
- Check service account has "Secret Manager Secret Accessor" role

### "Could not connect to Cloud SQL"

- Verify `--add-cloudsql-instances` flag is correct
- Check DATABASE_URL secret format
- Ensure Cloud SQL instance is running

### "OAuth error"

- Verify Google OAuth redirect URIs include your Cloud Run URL
- Check GOOGLE_CLIENT_ID and GOOGLE_CLIENT_SECRET secrets

### "502 Bad Gateway"

- Check logs for application errors
- Verify application starts within timeout
- Check memory limits (increase if needed)

### "Service Unavailable"

- Check if Cloud Run service is deployed
- Verify region is correct
- Check if application is listening on port 8080

---

## Updating the Deployment

### Code Changes

```bash
# GitHub Actions: Just push to main/uat branch

# Manual:
docker build -t gcr.io/$PROJECT_ID/fastinghero-uat:v2 .
docker push gcr.io/$PROJECT_ID/fastinghero-uat:v2
gcloud run deploy fastinghero-uat --image gcr.io/$PROJECT_ID/fastinghero-uat:v2 --region us-central1
```

### Environment Variables

```bash
gcloud run services update fastinghero-uat \
  --region us-central1 \
  --set-env-vars NEW_VAR=value
```

### Secrets

```bash
# Update secret in Secret Manager
echo -n "new-secret-value" | gcloud secrets versions add SECRET_NAME --data-file=-

# Restart service to pick up new version
gcloud run services update fastinghero-uat --region us-central1
```

---

## Rolling Back

```bash
# List revisions
gcloud run revisions list --service fastinghero-uat --region us-central1

# Roll back to specific revision
gcloud run services update-traffic fastinghero-uat \
  --region us-central1 \
  --to-revisions REVISION_NAME=100
```

---

## Cost Optimization

### Stay in Free Tier

- Keep min instances at 0 (scale to zero)
- Use smallest Cloud SQL instance (db-f1-micro)
- Delete old container images
- Monitor billing dashboard

### Reduce Costs

```bash
# Scale down Cloud SQL when not in use
gcloud sql instances patch fastinghero-uat-db --activation-policy=NEVER

# Scale back up
gcloud sql instances patch fastinghero-uat-db --activation-policy=ALWAYS
```

---

## Production Checklist

Before promoting to production:

- [ ] Test all OAuth flows
- [ ] Verify database migrations
- [ ] Load test application
- [ ] Set up monitoring alerts
- [ ] Configure custom domain
- [ ] Enable CDN (if serving static assets)
- [ ] Review security headers
- [ ] Test backup/restore procedures
- [ ] Document rollback procedures
- [ ] Set up error reporting (Sentry, etc.)

---

## Support

For issues:

1. Check Cloud Run logs
2. Check Cloud SQL logs
3. Review [Google Cloud Run documentation](https://cloud.google.com/run/docs)
4. Check GitHub Issues
