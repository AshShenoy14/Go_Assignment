import React, { useEffect, useState } from 'react';
import { useAuth } from '../contexts/AuthContext';
import { getAppointments } from '../services/api';
import Navbar from '../components/Navbar';
import { Clock, Calendar as CalendarIcon } from 'lucide-react';
import Calendar from '../components/Calendar';
import DoctorAvailability from '../components/DoctorAvailability';

const DoctorDashboard = () => {
    const { user } = useAuth();
    const [appointments, setAppointments] = useState([]);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        const fetchDoctorData = async () => {
            try {
                // Fetch appointments for this doctor
                const data = await getAppointments({ doctor_id: user.id });
                setAppointments(data || []);
            } catch (e) {
                console.error("Failed to fetch doctor data", e);
            } finally {
                setLoading(false);
            }
        };

        if (user?.id) {
            fetchDoctorData();
        }
    }, [user]);

    return (
        <div className="min-h-screen bg-zinc-950 text-white">
            <Navbar />
            <div className="max-w-7xl mx-auto px-4 py-8">
                <div className="flex items-center justify-between mb-8">
                    <div>
                        <h1 className="text-3xl font-bold">Doctor Dashboard</h1>
                        <p className="text-zinc-400 mt-1">Welcome, Dr. {user.username}</p>
                    </div>
                </div>

                <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
                    {/* Availability Settings - Sidebar on desktop */}
                    <div className="lg:col-span-1">
                        <DoctorAvailability />
                    </div>

                    {/* Appointments Calendar/List - Main Content */}
                    <div className="lg:col-span-2">
                        <div className="bg-zinc-900 p-6 rounded-xl border border-zinc-800 mb-6">
                            <div className="flex items-center justify-between mb-6">
                                <h2 className="text-xl font-semibold flex items-center">
                                    <CalendarIcon className="h-5 w-5 mr-2 text-purple-500" />
                                    My Schedule
                                </h2>
                            </div>

                            {loading ? (
                                <div className="flex justify-center py-10">
                                    <div className="animate-spin h-8 w-8 border-2 border-blue-500 border-t-transparent rounded-full"></div>
                                </div>
                            ) : appointments.length === 0 ? (
                                <div className="text-center py-10 text-zinc-500">
                                    No appointments scheduled properly yet.
                                </div>
                            ) : (
                                <Calendar
                                    events={appointments.map(app => ({
                                        id: app.id,
                                        title: app.title || 'Appointment',
                                        start: `${app.date}T${app.time}`,
                                        end: new Date(new Date(`${app.date}T${app.time}`).getTime() + (app.duration || 30) * 60000).toISOString(),
                                        backgroundColor: app.status === 'cancelled' ? '#EF4444' : '#3B82F6',
                                        borderColor: app.status === 'cancelled' ? '#EF4444' : '#3B82F6',
                                        extendedProps: app
                                    }))}
                                    onEventClick={(arg) => alert(`Appointment: ${arg.event.title}\nStatus: ${arg.event.extendedProps.status}`)}
                                />
                            )}
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );
};

export default DoctorDashboard;
