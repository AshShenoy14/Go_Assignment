import React, { useState, useEffect } from 'react';
import { motion } from 'framer-motion';
import { ChevronLeft, ChevronRight, Clock } from 'lucide-react';

const CalendarView = ({ appointments, onDateClick, onAppointmentClick }) => {
    const [currentDate, setCurrentDate] = useState(new Date());
    const [calendarDays, setCalendarDays] = useState([]);

    useEffect(() => {
        generateCalendar(currentDate);
    }, [currentDate, appointments]);

    const generateCalendar = (date) => {
        const year = date.getFullYear();
        const month = date.getMonth();

        const firstDayOfMonth = new Date(year, month, 1);
        const lastDayOfMonth = new Date(year, month + 1, 0);

        const daysInMonth = lastDayOfMonth.getDate();
        const startDayOfWeek = firstDayOfMonth.getDay(); // 0 = Sunday

        const days = [];

        // Previous month filler days
        for (let i = 0; i < startDayOfWeek; i++) {
            days.push({ day: null, date: null, type: 'empty' });
        }

        // Current month days
        for (let i = 1; i <= daysInMonth; i++) {
            const currentDayDate = new Date(year, month, i);
            const dateString = currentDayDate.toISOString().split('T')[0];

            const daysAppointments = appointments.filter(app => app.date === dateString);

            days.push({
                day: i,
                date: dateString,
                appointments: daysAppointments,
                type: 'day'
            });
        }

        setCalendarDays(days);
    };

    const nextMonth = () => {
        setCurrentDate(new Date(currentDate.getFullYear(), currentDate.getMonth() + 1, 1));
    };

    const prevMonth = () => {
        setCurrentDate(new Date(currentDate.getFullYear(), currentDate.getMonth() - 1, 1));
    };

    const monthNames = ["January", "February", "March", "April", "May", "June",
        "July", "August", "September", "October", "November", "December"
    ];

    return (
        <div className="bg-zinc-900 rounded-xl border border-zinc-800 overflow-hidden">
            {/* Header */}
            <div className="flex items-center justify-between p-6 border-b border-zinc-800">
                <h2 className="text-xl font-bold">
                    {monthNames[currentDate.getMonth()]} {currentDate.getFullYear()}
                </h2>
                <div className="flex gap-2">
                    <button onClick={prevMonth} className="p-2 hover:bg-zinc-800 rounded-lg text-zinc-400 hover:text-white transition-colors">
                        <ChevronLeft className="h-5 w-5" />
                    </button>
                    <button onClick={nextMonth} className="p-2 hover:bg-zinc-800 rounded-lg text-zinc-400 hover:text-white transition-colors">
                        <ChevronRight className="h-5 w-5" />
                    </button>
                </div>
            </div>

            {/* Days Header */}
            <div className="grid grid-cols-7 border-b border-zinc-800">
                {['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat'].map(day => (
                    <div key={day} className="py-3 text-center text-sm font-medium text-zinc-500">
                        {day}
                    </div>
                ))}
            </div>

            {/* Calendar Grid */}
            <div className="grid grid-cols-7 auto-rows-fr">
                {calendarDays.map((day, index) => (
                    <div
                        key={index}
                        className={`min-h-[120px] border-b border-r border-zinc-800/50 p-2 relative group transition-colors ${day.type === 'day' ? 'hover:bg-zinc-800/20 cursor-pointer' : 'bg-zinc-900/50'
                            }`}
                        onClick={() => day.type === 'day' && onDateClick(day.date)}
                    >
                        {day.type === 'day' && (
                            <>
                                <span className={`text-sm font-medium ${new Date().toISOString().split('T')[0] === day.date
                                        ? 'bg-blue-600 text-white w-6 h-6 rounded-full flex items-center justify-center'
                                        : 'text-zinc-400'
                                    }`}>
                                    {day.day}
                                </span>

                                <div className="mt-2 space-y-1">
                                    {day.appointments.map(app => (
                                        <motion.div
                                            key={app.id}
                                            initial={{ opacity: 0, scale: 0.9 }}
                                            animate={{ opacity: 1, scale: 1 }}
                                            className="text-xs bg-blue-500/10 text-blue-400 px-2 py-1 rounded border border-blue-500/20 truncate hover:bg-blue-500/20 hover:border-blue-500/40 transition-colors"
                                            onClick={(e) => {
                                                e.stopPropagation();
                                                onAppointmentClick(app);
                                            }}
                                        >
                                            <div className="flex items-center gap-1">
                                                <Clock className="w-3 h-3" />
                                                <span className="font-semibold">{app.time}</span>
                                            </div>
                                            <div className="truncate">{app.title}</div>
                                        </motion.div>
                                    ))}
                                </div>

                                {/* Add button placeholder on hover */}
                                <div className="absolute top-2 right-2 opacity-0 group-hover:opacity-100 transition-opacity">
                                    <div className="w-6 h-6 rounded-full bg-zinc-800 flex items-center justify-center text-zinc-400 hover:bg-blue-600 hover:text-white transition-colors">
                                        +
                                    </div>
                                </div>
                            </>
                        )}
                    </div>
                ))}
            </div>
        </div>
    );
};

export default CalendarView;
