import axios from 'axios';
import { Gadget, GadgetRequest, GadgetSession } from '../types';

const API_BASE_URL = import.meta.env.VITE_API_URL || '/api';

export const api = {
  async getGadgets(): Promise<Gadget[]> {
    const response = await axios.get(`${API_BASE_URL}/gadgets`);
    return response.data;
  },

  async getSessions(): Promise<GadgetSession[]> {
    const response = await axios.get(`${API_BASE_URL}/sessions`);
    return response.data;
  },

  async startSession(request: GadgetRequest): Promise<GadgetSession> {
    const response = await axios.post(`${API_BASE_URL}/sessions`, request);
    return response.data;
  },

  async stopSession(sessionId: string): Promise<void> {
    await axios.delete(`${API_BASE_URL}/sessions/${sessionId}`);
  },

  getWebSocketUrl(sessionId: string): string {
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const host = import.meta.env.VITE_WS_URL || window.location.host;
    return `${protocol}//${host}/ws/${sessionId}`;
  },
};
