import React, { useState, useEffect } from 'react';
import { useAuth } from '../contexts/AuthContext';
import { setAvailability, getAvailability } from '../services/api';
import Button from './Button';
import Input from './Input';
import { Clock, Check } from 'lucide-react';
import { motion } from 'framer-motion';

const DAYS = [
    { id: 1, name: 'Monday' },
    { id: 2, name: 'Tuesday' },
    { id: 3, name: 'Wednesday' },
    { id: 4, name: 'Thursday' },
    { id: 5, name: 'Friday' },
    { id: 6, name: 'Saturday' },
    { id: 0, name: 'Sunday' },
];

const DoctorAvailability = () => {
    const { user } = useAuth();
    const [schedule, setSchedule] = useState({});
    const [loading, setLoading] = useState(false);
    const [saving, setSaving] = useState(null); // day id being saved

    // Load initial availability
    useEffect(() => {
        const loadAvailability = async () => {
            if (!user) return;
            // For now, let's just default to empty or standard 09:00 - 17:00
            // In a real app we would fetch for all days.
            // Let's implement fetching one by one or just lazy load.
            // For UI simplicity, let's pre-fill with defaults if not found (logic handled by backend 404/null?)
            // We'll iterate and fetch.
            const newSchedule = {};
            for (let day of DAYS) {
                try {
                    const av = await getAvailability(user.id, day.id);
                    if (av) {
                        newSchedule[day.id] = { start: av.start_time, end: av.end_time };
                    }
                } catch (e) {
                    // Not set, ignore
                }
            }
            // Merge with defaults for UI if missing? No, keep empty to show "Not Set"
            // Actually better to have defaults for easy editing
            setSchedule(prev => ({ ...prev, ...newSchedule }));
        };
        loadAvailability();
    }, [user]);

    const handleSave = async (dayId) => {
        const daySchedule = schedule[dayId] || { start: '09:00', end: '17:00' };
        setSaving(dayId);
        try {
            await setAvailability({
                doctor_id: parseInt(user.id),
                day_of_week: dayId,
                start_time: daySchedule.start,
                end_time: daySchedule.end
            });
            // Show success ?
        } catch (error) {
            console.error("Failed to set availability", error);
            alert("Failed to save");
        } finally {
            setSaving(null);
        }
    };

    const updateTime = (dayId, field, value) => {
        setSchedule(prev => ({
            ...prev,
            [dayId]: {
                ...(prev[dayId] || { start: '09:00', end: '17:00' }),
                [field]: value
            }
        }));
    };

    return (
        <div className="bg-zinc-900 rounded-xl p-6 border border-zinc-800">
            <h2 className="text-xl font-bold text-white mb-4 flex items-center">
                <Clock className="w-5 h-5 mr-2 text-blue-500" />
                Weekly Availability
            </h2>
            <p className="text-zinc-400 text-sm mb-6">
                Set your working hours for each day. Patients will only be able to book slots within these times.
            </p>

            <div className="space-y-4">
                {DAYS.map(day => {
                    const daySch = schedule[day.id] || { start: '09:00', end: '17:00' };
                    return (
                        <div key={day.id} className="flex items-center gap-4 p-3 bg-zinc-950/50 rounded-lg border border-zinc-800/50">
                            <div className="w-24 font-medium text-zinc-300">{day.name}</div>
                            <div className="flex items-center gap-2">
                                <input
                                    type="time"
                                    value={daySch.start}
                                    onChange={(e) => updateTime(day.id, 'start', e.target.value)}
                                    className="bg-zinc-900 border border-zinc-700 rounded px-2 py-1 text-sm text-white focus:border-blue-500 outline-none"
                                />
                                <span className="text-zinc-500">-</span>
                                <input
                                    type="time"
                                    value={daySch.end}
                                    onChange={(e) => updateTime(day.id, 'end', e.target.value)}
                                    className="bg-zinc-900 border border-zinc-700 rounded px-2 py-1 text-sm text-white focus:border-blue-500 outline-none"
                                />
                            </div>
                            <button
                                onClick={() => handleSave(day.id)}
                                disabled={saving === day.id}
                                className={`ml-auto p-2 rounded-lg transition-colors ${saving === day.id
                                        ? 'bg-blue-500/20 text-blue-500'
                                        : 'bg-zinc-800 text-zinc-400 hover:bg-blue-600 hover:text-white'
                                    }`}
                            >
                                {saving === day.id ? <div className="animate-spin h-4 w-4 border-2 border-current border-t-transparent rounded-full" /> : <Check className="w-4 h-4" />}
                            </button>
                        </div>
                    );
                })}
            </div>
        </div>
    );
};

export default DoctorAvailability;
