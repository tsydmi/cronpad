import axios from "axios";

export async function initAxiosInterceptors(keycloak){
    console.log('axios interceptors init')
    await axios.interceptors.request.use(function (config) {
        console.log('axios interceptor')
        config.headers.Authorization = `Bearer ${keycloak.token}`

        return config;
    });
}