import Axios, { InternalAxiosRequestConfig, AxiosRequestConfig } from "axios";

import { useNotifications } from "@/components/ui/notifications";
import { env } from "@/config/env";
//import { paths } from "@/config/paths";

function authRequestInterceptor(config: InternalAxiosRequestConfig) {
  if (config.headers) {
    config.headers.Accept = "application/json";
  }

  config.withCredentials = true;
  return config;
}

export const axiosInstance = Axios.create({
  baseURL: env.API_URL,
});

axiosInstance.interceptors.request.use(authRequestInterceptor);
axiosInstance.interceptors.response.use(
  (response) => {
    return response.data;
  },
  (error) => {
    const message = error.response?.data?.message || error.message;
    useNotifications.getState().addNotification({
      type: "error",
      title: "Error",
      message,
    });

    //if (error.response?.status === 401) {
    //  const searchParams = new URLSearchParams();
    //  const redirectTo =
    //    searchParams.get("redirectTo") || window.location.pathname;
    //  window.location.href = paths.auth.login.getHref(redirectTo);
    //}

    return Promise.reject(error);
  },
);

type DataAxiosInstance = {
  get<T = unknown>(url: string, config?: AxiosRequestConfig): Promise<T>;
  post<T = unknown>(
    url: string,
    data?: any,
    config?: AxiosRequestConfig,
  ): Promise<T>;
  put<T = unknown>(
    url: string,
    data?: any,
    config?: AxiosRequestConfig,
  ): Promise<T>;
  patch<T = unknown>(
    url: string,
    data?: any,
    config?: AxiosRequestConfig,
  ): Promise<T>;
  delete<T = unknown>(url: string, config?: AxiosRequestConfig): Promise<T>;
};

export const api = axiosInstance as unknown as DataAxiosInstance;
