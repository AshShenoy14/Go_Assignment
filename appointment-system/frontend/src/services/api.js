import axios from 'axios';

const API_URL = 'http://localhost:8080/api';

const api = axios.create({
    baseURL: API_URL,
    headers: {
        'Content-Type': 'application/json',
    },
});

// Add a request interceptor to include the JWT token
api.interceptors.request.use(
    (config) => {
        const token = localStorage.getItem('token');
        if (token) {
            config.headers.Authorization = `Bearer ${token}`;
        }
        return config;
    },
    (error) => {
        return Promise.reject(error);
    }
);

export const login = async (email, password) => {
    const response = await api.post('/users/login', { email, password });
    return response.data;
};

export const register = async (username, email, password) => {
    const response = await api.post('/users/register', { username, email, password });
    return response.data;
};

export const updateUser = async (id, userData) => {
    const response = await api.put(`/users/${id}`, userData);
    return response.data;
};

export const getAppointments = async (criteria) => {
    // criteria can be userid (string/int) or object { user_id, doctor_id, date }
    let query = '';
    if (typeof criteria === 'object') {
        const params = new URLSearchParams();
        if (criteria.user_id) params.append('user_id', criteria.user_id);
        if (criteria.doctor_id) params.append('doctor_id', criteria.doctor_id);
        if (criteria.date) params.append('date', criteria.date);
        query = `?${params.toString()}`;
    } else if (criteria) {
        query = `?user_id=${criteria}`;
    }
    const response = await api.get(`/appointments${query}`);
    return response.data;
};

export const createAppointment = async (appointmentData) => {
    const response = await api.post('/appointments', appointmentData);
    return response.data;
};

export const updateAppointment = async (id, appointmentData) => {
    const response = await api.put(`/appointments/${id}`, appointmentData);
    return response.data;
};

export const deleteAppointment = async (id) => {
    const response = await api.delete(`/appointments/${id}`);
    return response.data;
};

export const getSlots = async (doctorId, date) => {
    const response = await api.get(`/appointments/slots?doctor_id=${doctorId}&date=${date}`);
    return response.data;
};

export const setAvailability = async (availabilityData) => {
    const response = await api.post(`/appointments/availability`, availabilityData);
    return response.data;
};

export const getAvailability = async (doctorId, dayOfWeek) => {
    const response = await api.get(`/appointments/availability?doctor_id=${doctorId}&day_of_week=${dayOfWeek}`);
    return response.data;
};

export const getDoctors = async () => {
    const response = await api.get('/users?role=doctor');
    return response.data;
};

export default api;

export const getStats = async () => {
    const response = await api.get('/stats'); // This might need to be /users/stats or /appointments/stats depending on what we want. 
    // We created /stats in both services. API Gateway routes /api/users/stats -> user-service and /api/appointments/stats -> appointment-service.
    // Let's assume we want a combined stats or specific. 
    // For Admin Dashboard, we probably want both. 
    // Let's add separate methods or one combined if gateway aggregated them (it doesn't).
    // Let's use /api/users/stats and /api/appointments/stats.
    const userStats = await api.get('/users/stats');
    const apptStats = await api.get('/appointments/stats');
    return { ...userStats.data, ...apptStats.data };
};

export const getNotifications = async (userId) => {
    const response = await api.get(`/notifications?user_id=${userId}`);
    return response.data;
};

export const cancelAppointment = async (id, reason) => {
    const response = await api.post(`/appointments/${id}/cancel`, { reason });
    return response.data;
};

export const getUsers = async (role) => {
    const query = role ? `?role=${role}` : '';
    const response = await api.get(`/users${query}`);
    return response.data;
};

export const deleteUser = async (id) => {
    const response = await api.delete(`/users/${id}`);
    return response.data;
};
