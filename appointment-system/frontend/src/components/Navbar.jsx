import React, { useState, useEffect } from 'react';
import { useAuth } from '../contexts/AuthContext';
import Button from './Button';
import { LogOut, Calendar, Bell } from 'lucide-react';
import { useNavigate, Link } from 'react-router-dom';
import { getNotifications } from '../services/api';
import { motion, AnimatePresence } from 'framer-motion';

const Navbar = () => {
    const { user, logout } = useAuth();
    const navigate = useNavigate();
    const [notifications, setNotifications] = useState([]);
    const [showNotifications, setShowNotifications] = useState(false);

    useEffect(() => {
        const fetchNotifications = async () => {
            try {
                if (user?.id) {
                    const data = await getNotifications(user.id);
                    setNotifications(data || []);
                }
            } catch (error) {
                console.error("Failed to fetch notifications", error);
            }
        };

        fetchNotifications();
        const interval = setInterval(fetchNotifications, 30000); // Poll every 30s
        return () => clearInterval(interval);
    }, [user]);

    const handleLogout = () => {
        logout();
        navigate('/login');
    };

    const unreadCount = notifications.filter(n => !n.is_read).length;

    return (
        <nav className="border-b border-zinc-800 bg-zinc-900/50 backdrop-blur-xl sticky top-0 z-50">
            <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
                <div className="flex items-center justify-between h-16">
                    <div className="flex items-center space-x-3">
                        <Link to="/dashboard" className="flex items-center space-x-3">
                            <div className="p-2 bg-blue-600/10 rounded-lg">
                                <Calendar className="h-6 w-6 text-blue-500" />
                            </div>
                            <span className="text-xl font-bold bg-gradient-to-r from-blue-400 to-purple-500 bg-clip-text text-transparent">
                                AppointmentSys
                            </span>
                        </Link>
                    </div>

                    <div className="flex items-center gap-4">
                        {/* Notifications */}
                        <div className="relative">
                            <button
                                onClick={() => setShowNotifications(!showNotifications)}
                                className="relative p-2 text-zinc-400 hover:text-white transition-colors"
                            >
                                <Bell className="h-5 w-5" />
                                {unreadCount > 0 && (
                                    <span className="absolute top-1 right-1 h-2 w-2 bg-red-500 rounded-full"></span>
                                )}
                            </button>

                            <AnimatePresence>
                                {showNotifications && (
                                    <motion.div
                                        initial={{ opacity: 0, y: 10 }}
                                        animate={{ opacity: 1, y: 0 }}
                                        exit={{ opacity: 0, y: 10 }}
                                        className="absolute right-0 mt-2 w-80 bg-zinc-900 border border-zinc-800 rounded-xl shadow-xl overflow-hidden z-50"
                                    >
                                        <div className="p-3 border-b border-zinc-800 font-semibold text-sm">Notifications</div>
                                        <div className="max-h-64 overflow-y-auto">
                                            {notifications.length === 0 ? (
                                                <div className="p-4 text-center text-zinc-500 text-sm">No notifications</div>
                                            ) : (
                                                notifications.map(n => (
                                                    <div key={n.id} className="p-3 border-b border-zinc-800 hover:bg-zinc-800/50 transition-colors text-sm">
                                                        <p className="text-zinc-300">{n.message}</p>
                                                        <p className="text-xs text-zinc-500 mt-1">{new Date(n.created_at).toLocaleString()}</p>
                                                    </div>
                                                ))
                                            )}
                                        </div>
                                    </motion.div>
                                )}
                            </AnimatePresence>
                        </div>

                        <span className="text-sm text-zinc-400 hidden sm:block">
                            {user?.email || user?.username}
                        </span>
                        <Button variant="ghost" onClick={() => navigate('/profile')} className="text-sm">
                            Profile
                        </Button>
                        <Button variant="ghost" onClick={handleLogout} className="text-sm">
                            <LogOut className="h-4 w-4 mr-2" />
                            Sign Out
                        </Button>
                    </div>
                </div>
            </div>
        </nav>
    );
};

export default Navbar;
