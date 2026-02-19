import React, { useEffect, useState } from 'react';
import { useAuth } from '../contexts/AuthContext';
import { getAppointments, createAppointment, getDoctors, cancelAppointment, getSlots } from '../services/api';
import Navbar from '../components/Navbar';
import Button from '../components/Button';
import Input from '../components/Input';
import { Search, Filter, Calendar as CalendarIcon, Clock, Trash2, X, User, ChevronRight } from 'lucide-react';
import { motion, AnimatePresence } from 'framer-motion';

const PatientDashboard = () => {
    const { user } = useAuth();
    const [appointments, setAppointments] = useState([]);
    const [doctors, setDoctors] = useState([]);
    const [loading, setLoading] = useState(true);
    const [isModalOpen, setIsModalOpen] = useState(false);
    const [viewMode, setViewMode] = useState('appointments'); // 'appointments' or 'find_doctors'
    const [searchQuery, setSearchQuery] = useState('');
    const [specializationFilter, setSpecializationFilter] = useState('');
    const [selectedDoctor, setSelectedDoctor] = useState(null);
    const [slots, setSlots] = useState([]);
    const [loadingSlots, setLoadingSlots] = useState(false);

    const [formData, setFormData] = useState({
        title: '',
        description: '',
        date: '',
        time: '',
        duration: 30
    });

    const fetchData = async () => {
        try {
            const [apptData, doctorsData] = await Promise.all([
                getAppointments({ user_id: user.id }),
                getDoctors()
            ]);
            setAppointments(apptData || []);
            setDoctors(doctorsData || []);
        } catch (error) {
            console.error("Failed to fetch patient data", error);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        if (user?.id) fetchData();
    }, [user]);

    const handleBookClick = (doctor) => {
        setSelectedDoctor(doctor);
        setFormData({
            title: '',
            description: '',
            date: '', // Reset date
            time: '',
            duration: 30
        });
        setSlots([]);
        setIsModalOpen(true);
    };

    const handleDateChange = async (e) => {
        const date = e.target.value;
        setFormData({ ...formData, date, time: '' });
        setSlots([]); // clear old slots

        if (date && selectedDoctor) {
            setLoadingSlots(true);
            try {
                const availableSlots = await getSlots(selectedDoctor.id, date);
                setSlots(availableSlots || []);
            } catch (error) {
                console.error("Failed to fetch slots", error);
            } finally {
                setLoadingSlots(false);
            }
        }
    };

    const handleSubmit = async (e) => {
        e.preventDefault();
        try {
            await createAppointment({
                ...formData,
                user_id: user.id,
                doctor_id: selectedDoctor.id,
                duration: parseInt(formData.duration)
            });
            setIsModalOpen(false);
            fetchData();
            setViewMode('appointments'); // Switch back to list
            alert("Appointment booked successfully!");
        } catch (error) {
            alert("Failed to book appointment: " + (error.response?.data?.error || error.message));
        }
    };

    const handleCancel = async (id) => {
        const reason = prompt("Please enter a reason for cancellation:");
        if (reason === null) return;

        try {
            await cancelAppointment(id, reason || "User cancelled");
            fetchData();
        } catch (error) {
            alert("Failed to cancel appointment");
        }
    };

    // Filter Logic
    const filteredAppointments = appointments.filter(app =>
        app.title.toLowerCase().includes(searchQuery.toLowerCase())
    );

    const filteredDoctors = doctors.filter(doc => {
        const matchesName = doc.username.toLowerCase().includes(searchQuery.toLowerCase());
        const matchesSpec = specializationFilter ? doc.specialization === specializationFilter : true;
        return matchesName && matchesSpec;
    });

    const specializations = [...new Set(doctors.map(d => d.specialization).filter(Boolean))];

    return (
        <div className="min-h-screen bg-zinc-950 text-white">
            <Navbar />
            <main className="max-w-7xl mx-auto px-4 py-8">
                <div className="flex flex-col md:flex-row justify-between items-center mb-8 gap-4">
                    <div>
                        <h1 className="text-3xl font-bold">Patient Dashboard</h1>
                        <p className="text-zinc-400">Manage your health and appointments.</p>
                    </div>
                    <div className="flex gap-2">
                        <Button
                            variant={viewMode === 'appointments' ? 'primary' : 'secondary'}
                            onClick={() => setViewMode('appointments')}
                        >
                            My Appointments
                        </Button>
                        <Button
                            variant={viewMode === 'find_doctors' ? 'primary' : 'secondary'}
                            onClick={() => setViewMode('find_doctors')}
                        >
                            Find a Doctor
                        </Button>
                    </div>
                </div>

                {viewMode === 'appointments' && (
                    <div className="space-y-6">
                        <div className="relative max-w-md">
                            <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-zinc-500" />
                            <input
                                type="text"
                                placeholder="Search appointments..."
                                value={searchQuery}
                                onChange={(e) => setSearchQuery(e.target.value)}
                                className="w-full bg-zinc-900 border border-zinc-800 rounded-lg pl-10 pr-4 py-2 text-sm focus:outline-none focus:border-blue-500"
                            />
                        </div>

                        {loading ? <div className="text-center py-10">Loading...</div> : (
                            <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
                                {filteredAppointments.map(app => (
                                    <div key={app.id} className="bg-zinc-900 rounded-xl p-6 border border-zinc-800 relative group">
                                        <div className="absolute top-4 right-4 opacity-0 group-hover:opacity-100 transition-opacity">
                                            {app.status !== 'cancelled' && (
                                                <button onClick={() => handleCancel(app.id)} className="text-zinc-500 hover:text-red-500" title="Cancel Appointment">
                                                    <Trash2 className="h-5 w-5" />
                                                </button>
                                            )}
                                        </div>
                                        <div className="flex items-start justify-between mb-2">
                                            <div className="p-2 bg-blue-500/10 rounded-lg">
                                                <CalendarIcon className="h-5 w-5 text-blue-400" />
                                            </div>
                                            <span className={`text-xs px-2 py-1 rounded border capitalize ${app.status === 'cancelled' ? 'bg-red-500/10 text-red-400 border-red-500/20' : 'bg-green-500/10 text-green-400 border-green-500/20'
                                                }`}>
                                                {app.status}
                                            </span>
                                        </div>
                                        <h3 className="text-lg font-semibold mb-1">{app.title}</h3>
                                        <p className="text-zinc-400 text-sm mb-4 line-clamp-2">{app.description}</p>

                                        <div className="space-y-2 text-sm text-zinc-500">
                                            <div className="flex items-center"><CalendarIcon className="h-4 w-4 mr-2" /> {app.date}</div>
                                            <div className="flex items-center"><Clock className="h-4 w-4 mr-2" /> {app.time} ({app.duration}m)</div>
                                            {app.doctor_id && (
                                                <div className="flex items-center text-blue-400">
                                                    <User className="h-4 w-4 mr-2" />
                                                    {doctors.find(d => d.id === app.doctor_id)?.username || `Doctor ID: ${app.doctor_id}`}
                                                </div>
                                            )}
                                        </div>
                                    </div>
                                ))}
                                {filteredAppointments.length === 0 && (
                                    <div className="col-span-full text-center py-10 text-zinc-500">No appointments found.</div>
                                )}
                            </div>
                        )}
                    </div>
                )}

                {viewMode === 'find_doctors' && (
                    <div className="space-y-6">
                        <div className="flex gap-4">
                            <div className="relative flex-1">
                                <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-zinc-500" />
                                <input
                                    type="text"
                                    placeholder="Search doctors by name..."
                                    value={searchQuery}
                                    onChange={(e) => setSearchQuery(e.target.value)}
                                    className="w-full bg-zinc-900 border border-zinc-800 rounded-lg pl-10 pr-4 py-2 text-sm focus:outline-none focus:border-blue-500"
                                />
                            </div>
                            <select
                                value={specializationFilter}
                                onChange={(e) => setSpecializationFilter(e.target.value)}
                                className="bg-zinc-900 border border-zinc-800 rounded-lg px-4 py-2 text-sm focus:outline-none focus:border-blue-500 min-w-[200px]"
                            >
                                <option value="">All Specializations</option>
                                {specializations.map(spec => (
                                    <option key={spec} value={spec}>{spec}</option>
                                ))}
                            </select>
                        </div>

                        <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
                            {filteredDoctors.map(doc => (
                                <div key={doc.id} className="bg-zinc-900 rounded-xl p-6 border border-zinc-800 flex flex-col">
                                    <div className="flex items-center gap-4 mb-4">
                                        <div className="h-12 w-12 rounded-full bg-zinc-800 flex items-center justify-center text-xl font-bold text-zinc-500">
                                            {doc.username.charAt(0).toUpperCase()}
                                        </div>
                                        <div>
                                            <h3 className="text-lg font-semibold">Dr. {doc.username}</h3>
                                            <p className="text-blue-400 text-sm">{doc.specialization || 'General Physician'}</p>
                                        </div>
                                    </div>

                                    <div className="mt-auto pt-4 border-t border-zinc-800">
                                        <Button className="w-full" onClick={() => handleBookClick(doc)}>
                                            Book Appointment <ChevronRight className="h-4 w-4 ml-1" />
                                        </Button>
                                    </div>
                                </div>
                            ))}
                        </div>
                    </div>
                )}
            </main>

            {/* Booking Modal */}
            <AnimatePresence>
                {isModalOpen && selectedDoctor && (
                    <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/60 backdrop-blur-sm">
                        <motion.div
                            initial={{ opacity: 0, scale: 0.95 }}
                            animate={{ opacity: 1, scale: 1 }}
                            exit={{ opacity: 0, scale: 0.95 }}
                            className="bg-zinc-900 rounded-xl p-6 w-full max-w-lg border border-zinc-800 shadow-2xl max-h-[90vh] overflow-y-auto"
                        >
                            <div className="flex justify-between items-center mb-6 border-b border-zinc-800 pb-4">
                                <div>
                                    <h2 className="text-xl font-bold">Book with Dr. {selectedDoctor.username}</h2>
                                    <p className="text-zinc-400 text-sm">{selectedDoctor.specialization}</p>
                                </div>
                                <button onClick={() => setIsModalOpen(false)}><X className="text-zinc-500 hover:text-white" /></button>
                            </div>

                            <form onSubmit={handleSubmit} className="space-y-5">
                                <Input label="Reason for Visit (Title)" value={formData.title} onChange={(e) => setFormData({ ...formData, title: e.target.value })} required placeholder="e.g., Annual Checkup" />
                                <Input label="Description (Optional)" value={formData.description} onChange={(e) => setFormData({ ...formData, description: e.target.value })} />

                                <div>
                                    <label className="block text-sm font-medium text-zinc-400 mb-1">Date</label>
                                    <input
                                        type="date"
                                        className="w-full bg-zinc-950 border border-zinc-700 rounded-lg px-3 py-2 text-white focus:outline-none focus:border-blue-500"
                                        value={formData.date}
                                        onChange={handleDateChange}
                                        min={new Date().toISOString().split('T')[0]}
                                        required
                                    />
                                </div>

                                {formData.date && (
                                    <div>
                                        <label className="block text-sm font-medium text-zinc-400 mb-2">Available Slots</label>
                                        {loadingSlots ? (
                                            <div className="text-sm text-zinc-500 animate-pulse">Checking availability...</div>
                                        ) : slots.length > 0 ? (
                                            <div className="grid grid-cols-3 sm:grid-cols-4 gap-2">
                                                {slots.map(slot => (
                                                    <button
                                                        key={slot}
                                                        type="button"
                                                        onClick={() => setFormData({ ...formData, time: slot })}
                                                        className={`px-3 py-2 rounded-lg text-sm border transition-colors ${formData.time === slot
                                                                ? 'bg-blue-600 border-blue-600 text-white'
                                                                : 'bg-zinc-800 border-zinc-700 text-zinc-300 hover:bg-zinc-700'
                                                            }`}
                                                    >
                                                        {slot}
                                                    </button>
                                                ))}
                                            </div>
                                        ) : (
                                            <div className="text-sm text-red-400 bg-red-500/10 p-3 rounded border border-red-500/20">
                                                No slots available for this date.
                                            </div>
                                        )}
                                    </div>
                                )}

                                <Button type="submit" className="w-full" disabled={!formData.time || !formData.date}>Confirm Booking</Button>
                            </form>
                        </motion.div>
                    </div>
                )}
            </AnimatePresence>
        </div>
    );
};

export default PatientDashboard;
