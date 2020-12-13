import axios from "axios";

export async function initAxiosInterceptors(keycloak) {
    axios.interceptors.request.use(function (config) {
        config.headers.Authorization = `Bearer ${keycloak.token}`
        return config;
    });
}