import React, { createContext, useContext, useState, useEffect } from 'react';
import * as api from '../services/api';
import { jwtDecode } from "jwt-decode";

const AuthContext = createContext();

export const useAuth = () => useContext(AuthContext);

export const AuthProvider = ({ children }) => {
    const [user, setUser] = useState(null);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        // Check local storage for token
        const token = localStorage.getItem('token');
        if (token) {
            try {
                const decoded = jwtDecode(token);
                // Check if token is expired
                const currentTime = Date.now() / 1000;
                if (decoded.exp < currentTime) {
                    logout();
                } else {
                    setUser({
                        id: decoded.user_id,
                        email: decoded.email,
                        role: decoded.role,
                        token: token
                    });
                }
            } catch (e) {
                console.error("Failed to decode token", e);
                logout();
            }
        }
        setLoading(false);
    }, []);

    const login = async (email, password) => {
        try {
            const data = await api.login(email, password);
            // data contains: { user_id, username, email, role, token }

            // Store token
            localStorage.setItem('token', data.token);

            // Set user state
            setUser({
                id: data.user_id,
                email: data.email,
                username: data.username,
                role: data.role,
                token: data.token
            });

            return data;
        } catch (error) {
            throw error;
        }
    };

    const register = async (username, email, password, role) => {
        try {
            // Register now accepts role (optional) but we'll deal with UI later
            // For now, if we want to support role in register, we need to pass it
            const data = await api.register(username, email, password, role);

            // If registration logs the user in automatically (it returns the same structure as login)
            if (data.token) {
                localStorage.setItem('token', data.token);
                setUser({
                    id: data.user_id,
                    email: data.email,
                    username: data.username,
                    role: data.role,
                    token: data.token
                });
            }
            return data;
        } catch (error) {
            throw error;
        }
    };

    const logout = () => {
        setUser(null);
        localStorage.removeItem('token');
    };

    const hasRole = (role) => {
        return user?.role === role;
    };

    return (
        <AuthContext.Provider value={{ user, login, register, logout, loading, hasRole }}>
            {!loading && children}
        </AuthContext.Provider>
    );
};
