import React, { useState } from 'react';
import { useAuth } from '../contexts/AuthContext';
import Navbar from '../components/Navbar';
import Button from '../components/Button';
import Input from '../components/Input';
import { updateUser } from '../services/api'; // Make sure this is exported from api.js
import { User, Mail, Lock, Check, AlertCircle } from 'lucide-react';
import { motion, AnimatePresence } from 'framer-motion';

const Profile = () => {
    const { user, login } = useAuth(); // We might need a way to update the user in context/storage
    const [formData, setFormData] = useState({
        username: user?.username || '',
        email: user?.email || '',
        phone: user?.phone || '',
        specialization: user?.specialization || '',
        password: '',
        confirmPassword: ''
    });
    const [loading, setLoading] = useState(false);
    const [message, setMessage] = useState({ type: '', text: '' });

    const handleSubmit = async (e) => {
        e.preventDefault();
        setMessage({ type: '', text: '' });

        if (formData.password && formData.password !== formData.confirmPassword) {
            setMessage({ type: 'error', text: 'Passwords do not match' });
            return;
        }

        setLoading(true);
        try {
            const userId = user.id || user.user_id; // Handle both cases just in case
            const payload = {
                username: formData.username,
                email: formData.email,
                phone: formData.phone,
                specialization: formData.specialization,
                password: formData.password || undefined // Only send if not empty
            };

            const updatedUser = await updateUser(userId, payload);

            setMessage({ type: 'success', text: 'Profile updated successfully' });
            // update local storage if needed, but for now we rely on next login or refresh?
            // Ideally AuthContext exposed a verify/refresh method, or we manually update it.
            // Let's manually update localStorage for persistence consistency
            const storedUser = JSON.parse(localStorage.getItem('user'));
            const newUserData = { ...storedUser, ...updatedUser };
            localStorage.setItem('user', JSON.stringify(newUserData));

            // clear password fields
            setFormData(prev => ({ ...prev, password: '', confirmPassword: '' }));
        } catch (error) {
            setMessage({ type: 'error', text: error.response?.data?.error || 'Failed to update profile' });
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="min-h-screen bg-zinc-950 text-white">
            <Navbar />

            <main className="max-w-3xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
                <motion.div
                    initial={{ opacity: 0, y: 20 }}
                    animate={{ opacity: 1, y: 0 }}
                    className="bg-zinc-900 rounded-2xl border border-zinc-800 p-8 shadow-xl"
                >
                    <div className="mb-8">
                        <h1 className="text-3xl font-bold mb-2">Profile Settings</h1>
                        <p className="text-zinc-400">Manage your account information and preferences</p>
                    </div>

                    <AnimatePresence>
                        {message.text && (
                            <motion.div
                                initial={{ opacity: 0, height: 0 }}
                                animate={{ opacity: 1, height: 'auto' }}
                                exit={{ opacity: 0, height: 0 }}
                                className={`mb-6 p-4 rounded-lg flex items-center ${message.type === 'success' ? 'bg-green-500/10 text-green-500' : 'bg-red-500/10 text-red-500'
                                    }`}
                            >
                                {message.type === 'success' ? <Check className="h-5 w-5 mr-2" /> : <AlertCircle className="h-5 w-5 mr-2" />}
                                {message.text}
                            </motion.div>
                        )}
                    </AnimatePresence>

                    <form onSubmit={handleSubmit} className="space-y-6">
                        <div className="space-y-4">
                            <h2 className="text-xl font-semibold flex items-center">
                                <User className="h-5 w-5 mr-2 text-blue-500" />
                                Personal Information
                            </h2>
                            <div className="grid gap-4 md:grid-cols-2">
                                <Input
                                    label="Username"
                                    value={formData.username}
                                    onChange={(e) => setFormData({ ...formData, username: e.target.value })}
                                />
                                <Input
                                    label="Email"
                                    type="email"
                                    value={formData.email}
                                    onChange={(e) => setFormData({ ...formData, email: e.target.value })}
                                />
                                <Input
                                    label="Phone Number"
                                    type="tel"
                                    value={formData.phone}
                                    onChange={(e) => setFormData({ ...formData, phone: e.target.value })}
                                />
                                {user?.role === 'doctor' && (
                                    <div className="space-y-1">
                                        <label className="block text-sm font-medium text-zinc-400">Specialization</label>
                                        <select
                                            value={formData.specialization}
                                            onChange={(e) => setFormData({ ...formData, specialization: e.target.value })}
                                            className="w-full bg-zinc-950 border border-zinc-700 rounded-lg px-4 py-2 text-white focus:outline-none focus:border-blue-500 transition-colors"
                                        >
                                            <option value="">Select Specialization</option>
                                            <option value="General Physician">General Physician</option>
                                            <option value="Cardiologist">Cardiologist</option>
                                            <option value="Dermatologist">Dermatologist</option>
                                            <option value="Pediatrician">Pediatrician</option>
                                            <option value="Neurologist">Neurologist</option>
                                            <option value="Orthopedic">Orthopedic</option>
                                        </select>
                                    </div>
                                )}
                            </div>
                        </div>

                        <div className="pt-6 border-t border-zinc-800 space-y-4">
                            <h2 className="text-xl font-semibold flex items-center">
                                <Lock className="h-5 w-5 mr-2 text-blue-500" />
                                Security
                            </h2>
                            <p className="text-sm text-zinc-500">Leave blank to keep current password</p>
                            <div className="grid gap-4 md:grid-cols-2">
                                <Input
                                    label="New Password"
                                    type="password"
                                    value={formData.password}
                                    onChange={(e) => setFormData({ ...formData, password: e.target.value })}
                                />
                                <Input
                                    label="Confirm New Password"
                                    type="password"
                                    value={formData.confirmPassword}
                                    onChange={(e) => setFormData({ ...formData, confirmPassword: e.target.value })}
                                />
                            </div>
                        </div>

                        <div className="pt-6 border-t border-zinc-800 flex justify-end">
                            <Button type="submit" isLoading={loading} className="w-full sm:w-auto">
                                Save Changes
                            </Button>
                        </div>
                    </form>
                </motion.div>
            </main>
        </div>
    );
};

export default Profile;
