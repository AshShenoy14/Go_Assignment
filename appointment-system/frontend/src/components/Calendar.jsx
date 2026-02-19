import React from 'react';
import FullCalendar from '@fullcalendar/react';
import dayGridPlugin from '@fullcalendar/daygrid';
import timeGridPlugin from '@fullcalendar/timegrid';
import interactionPlugin from '@fullcalendar/interaction';

const Calendar = ({ events, onDateClick, onEventClick }) => {
    return (
        <div className='p-4 bg-zinc-900 rounded-xl border border-zinc-800 shadow-lg text-white'>
            <style>{`
        .fc {
          --fc-border-color: #27272a;
          --fc-daygrid-event-dot-width: 8px;
          --fc-list-event-dot-width: 8px;
          --fc-event-bg-color: #3b82f6;
          --fc-event-border-color: #2563eb;
          --fc-now-indicator-color: #f43f5e;
          --fc-today-bg-color: rgba(59, 130, 246, 0.1);
        }
        .fc-theme-standard td, .fc-theme-standard th {
          border-color: var(--fc-border-color);
        }
        .fc-col-header-cell-cushion, .fc-daygrid-day-number {
          color: #e4e4e7; /* zinc-200 */
        }
        .fc-event {
          cursor: pointer;
        }
      `}</style>
            <FullCalendar
                plugins={[dayGridPlugin, timeGridPlugin, interactionPlugin]}
                initialView="dayGridMonth"
                headerToolbar={{
                    left: 'prev,next today',
                    center: 'title',
                    right: 'dayGridMonth,timeGridWeek,timeGridDay'
                }}
                events={events} // Array of { title, start, end, id }
                dateClick={onDateClick}
                eventClick={onEventClick}
                height="auto"
                selectable={true}
                editable={false} // True if you want drag/drop (requires update handler)
                dayMaxEvents={true}
            />
        </div>
    );
};

export default Calendar;
