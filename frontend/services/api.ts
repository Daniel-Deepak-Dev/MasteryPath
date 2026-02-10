import axios from 'axios';

const BASE_API_URL = process.env.EXPO_PUBLIC_API_URL || 'http://localhost:8080/api';
console.log('🚀 API Configured URL:', BASE_API_URL);

const api = axios.create({
  baseURL: BASE_API_URL,
  timeout: 10000,
});

// ──── Types ────

export interface DashboardSubSkill {
  id: string;
  name: string;
  mastery: number;
}

export interface DashboardSkill {
  id: string;
  name: string;
  category: string;
  description: string;
  sub_skills: DashboardSubSkill[];
  overall_mastery: number;
}

// ──── API Functions ────

export const fetchDashboard = async (): Promise<DashboardSkill[]> => {
  const response = await api.get<DashboardSkill[]>('/skills/dashboard');
  return response.data;
};

export default api;
