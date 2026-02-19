import React from 'react';
import { useAuth } from '../contexts/AuthContext';
import AdminDashboard from '../components/AdminDashboard';
import DoctorDashboard from '../components/DoctorDashboard';
import PatientDashboard from '../components/PatientDashboard';
import Navbar from '../components/Navbar';

const Dashboard = () => {
    const { user, loading } = useAuth();

    if (loading) {
        return <div className="min-h-screen bg-zinc-950 flex items-center justify-center text-white">Loading...</div>;
    }

    if (!user) {
        return <div className="min-h-screen bg-zinc-950 flex items-center justify-center text-white">Please log in.</div>;
    }

    // Role-based rendering
    switch (user.role) {
        case 'admin':
            return <AdminDashboard />;
        case 'doctor':
            return <DoctorDashboard />;
        case 'patient':
            return <PatientDashboard />;
        default:
            // Fallback for unknown roles or if role is missing (default to patient view usually, or error)
            return <PatientDashboard />;
    }
};

export default Dashboard;
