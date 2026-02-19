import React, { useEffect, useState } from 'react';
import { useAuth } from '../contexts/AuthContext';
import { getStats, getUsers, deleteUser } from '../services/api';
import { Users, Calendar, UserX, Trash2 } from 'lucide-react';
import Navbar from '../components/Navbar';

const AdminDashboard = () => {
    const { user } = useAuth();
    const [stats, setStats] = useState({
        total_users: 0,
        total_doctors: 0,
        total_appointments: 0,
        cancelled_appointments: 0,
        completed_appointments: 0
    });
    const [users, setUsers] = useState([]);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        const fetchData = async () => {
            try {
                const statsData = await getStats();
                setStats(statsData);

                const usersData = await getUsers(); // Get all users
                setUsers(usersData || []);
            } catch (error) {
                console.error("Failed to fetch admin data", error);
            } finally {
                setLoading(false);
            }
        };

        fetchData();
    }, []);

    const handleDeleteUser = async (id) => {
        if (!window.confirm("Are you sure you want to delete this user?")) return;
        try {
            await deleteUser(id);
            setUsers(users.filter(u => u.id !== id));
        } catch (error) {
            alert("Failed to delete user");
        }
    };

    if (loading) return <div className="text-white text-center mt-20">Loading admin dashboard...</div>;

    return (
        <div className="min-h-screen bg-zinc-950 text-white">
            <Navbar />
            <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
                <h1 className="text-3xl font-bold mb-8">Admin Dashboard</h1>

                {/* Stats Cards */}
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
                    <div className="bg-zinc-900 p-6 rounded-xl border border-zinc-800">
                        <div className="flex items-center justify-between">
                            <div>
                                <p className="text-zinc-400 text-sm">Total Users</p>
                                <p className="text-2xl font-bold">{stats.total_users}</p>
                            </div>
                            <Users className="h-8 w-8 text-blue-500" />
                        </div>
                    </div>
                    <div className="bg-zinc-900 p-6 rounded-xl border border-zinc-800">
                        <div className="flex items-center justify-between">
                            <div>
                                <p className="text-zinc-400 text-sm">Total Doctors</p>
                                <p className="text-2xl font-bold">{stats.total_doctors}</p>
                            </div>
                            <UserX className="h-8 w-8 text-green-500" />
                        </div>
                    </div>
                    <div className="bg-zinc-900 p-6 rounded-xl border border-zinc-800">
                        <div className="flex items-center justify-between">
                            <div>
                                <p className="text-zinc-400 text-sm">Total Appointments</p>
                                <p className="text-2xl font-bold">{stats.total_appointments}</p>
                            </div>
                            <Calendar className="h-8 w-8 text-purple-500" />
                        </div>
                    </div>
                    <div className="bg-zinc-900 p-6 rounded-xl border border-zinc-800">
                        <div className="flex items-center justify-between">
                            <div>
                                <p className="text-zinc-400 text-sm">Cancelled</p>
                                <p className="text-2xl font-bold">{stats.cancelled_appointments}</p>
                            </div>
                            <UserX className="h-8 w-8 text-red-500" />
                        </div>
                    </div>
                </div>

                {/* Users List */}
                <div className="bg-zinc-900 rounded-xl border border-zinc-800 overflow-hidden">
                    <div className="px-6 py-4 border-b border-zinc-800">
                        <h2 className="text-xl font-semibold">User Management</h2>
                    </div>
                    <div className="overflow-x-auto">
                        <table className="w-full text-left">
                            <thead className="bg-zinc-900/50 text-zinc-400 text-sm">
                                <tr>
                                    <th className="px-6 py-3">ID</th>
                                    <th className="px-6 py-3">Username</th>
                                    <th className="px-6 py-3">Email</th>
                                    <th className="px-6 py-3">Role</th>
                                    <th className="px-6 py-3">Actions</th>
                                </tr>
                            </thead>
                            <tbody className="divide-y divide-zinc-800">
                                {users.map(u => (
                                    <tr key={u.id} className="hover:bg-zinc-800/50 transition-colors">
                                        <td className="px-6 py-4">{u.id}</td>
                                        <td className="px-6 py-4">{u.username}</td>
                                        <td className="px-6 py-4">{u.email}</td>
                                        <td className="px-6 py-4 capitalize">{u.role}</td>
                                        <td className="px-6 py-4">
                                            <button
                                                onClick={() => handleDeleteUser(u.id)}
                                                className="text-red-500 hover:text-red-400 transition-colors"
                                                title="Delete User"
                                            >
                                                <Trash2 className="h-5 w-5" />
                                            </button>
                                        </td>
                                    </tr>
                                ))}
                            </tbody>
                        </table>
                    </div>
                </div>
            </main>
        </div>
    );
};

export default AdminDashboard;
