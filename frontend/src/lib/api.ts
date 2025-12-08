import axios, { type AxiosInstance, type AxiosRequestConfig } from "axios";
import config from "@app/config.json";

export interface ApiResponse<T> {
  success: boolean;
  data: T;
  error: object;
  messages: string[];
}

// Custom axios instance that unwraps ApiResponse
interface CustomAxiosInstance
  extends Omit<AxiosInstance, "get" | "post" | "put" | "delete" | "patch"> {
  get<T>(url: string, config?: AxiosRequestConfig): Promise<T>;
  post<T>(url: string, data?: unknown, config?: AxiosRequestConfig): Promise<T>;
  put<T>(url: string, data?: unknown, config?: AxiosRequestConfig): Promise<T>;
  delete<T>(url: string, config?: AxiosRequestConfig): Promise<T>;
  patch<T>(
    url: string,
    data?: unknown,
    config?: AxiosRequestConfig
  ): Promise<T>;
}

const axiosInstance = axios.create({
  baseURL: `http://localhost:${config.rest_port}`,
  withCredentials: true,
});

axiosInstance.interceptors.response.use(
  (res) => res.data.data,
  (err) => {
    const errorMessage =
      err.response?.data?.error ||
      err.response?.data?.message ||
      err.message ||
      "An unexpected error occurred";
    return Promise.reject(new Error(errorMessage));
  }
);

const api = axiosInstance as unknown as CustomAxiosInstance;

export default api;
