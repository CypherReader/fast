import { useNavigate } from 'react-router-dom';
import { LogOut, User, Settings2 } from 'lucide-react';
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuLabel,
    DropdownMenuSeparator,
    DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import { useUser } from '@/hooks/use-user';

const UserMenu = () => {
    const navigate = useNavigate();
    const { user } = useUser();

    const handleLogout = () => {
        // Clear authentication token
        localStorage.removeItem('token');
        sessionStorage.removeItem('token');

        // Navigate to login page
        navigate('/login');
    };

    return (
        <DropdownMenu>
            <DropdownMenuTrigger asChild>
                <button className="w-8 h-8 rounded-full bg-primary/20 flex items-center justify-center hover:bg-primary/30 transition-colors focus:outline-none focus:ring-2 focus:ring-primary focus:ring-offset-2 focus:ring-offset-background">
                    <span className="text-sm font-medium text-primary">
                        {user?.name ? user.name.charAt(0).toUpperCase() : 'U'}
                    </span>
                </button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end" className="w-48">
                <DropdownMenuLabel>
                    <div className="flex flex-col space-y-1">
                        <p className="text-sm font-medium text-foreground">
                            {user?.name || 'User'}
                        </p>
                        <p className="text-xs text-muted-foreground truncate">
                            {user?.email || 'user@example.com'}
                        </p>
                    </div>
                </DropdownMenuLabel>
                <DropdownMenuSeparator />
                <DropdownMenuItem onClick={() => navigate('/profile')} className="cursor-pointer">
                    <User className="mr-2 h-4 w-4" />
                    <span>Profile</span>
                </DropdownMenuItem>
                <DropdownMenuItem onClick={() => navigate('/settings')} className="cursor-pointer">
                    <Settings2 className="mr-2 h-4 w-4" />
                    <span>Settings</span>
                </DropdownMenuItem>
                <DropdownMenuSeparator />
                <DropdownMenuItem onClick={handleLogout} className="cursor-pointer text-destructive focus:text-destructive">
                    <LogOut className="mr-2 h-4 w-4" />
                    <span>Logout</span>
                </DropdownMenuItem>
            </DropdownMenuContent>
        </DropdownMenu>
    );
};

export default UserMenu;
