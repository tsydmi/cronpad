import Vue from 'vue'
import App from './App.vue'
import vuetify from './plugins/vuetify'
import axios from 'axios'
import Keycloak from 'keycloak-js'
import './sass/variables.scss'

import router from '@/config/router'
import {initAxiosInterceptors} from '@/config/axios'
import dayjs from 'dayjs'

Vue.config.productionTip = false

const dayjsUtcPlugin = require('dayjs/plugin/utc')
dayjs.extend(dayjsUtcPlugin)

const getEnvironmentConfig = async () => {
    const config = await fetch('/environment.json')
    return await config.json()
}

getEnvironmentConfig().then(function (envJson){
    Vue.prototype.$http = axios
    axios.defaults.baseURL = envJson.VUE_APP_BACKEND_URL

    let initOptions = {
        url: envJson.VUE_APP_KEYCLOAK_URL, realm: 'cronpad', clientId: 'vue-frontend', onLoad: 'login-required'
    }

    let keycloak = Keycloak(initOptions);

    keycloak.init({onLoad: initOptions.onLoad})
        .then((auth) => {
            if (!auth) {
                window.location.reload();
            } else {
                console.log("Authenticated");
                initAxiosInterceptors(keycloak)

                new Vue({
                    el: '#app',
                    vuetify,
                    router,
                    render: h => h(App, {props: {keycloak: keycloak}})
                })
            }

            //Token Refresh
            setInterval(() => {
                keycloak.updateToken(70).catch(() => {
                    console.log('Failed to refresh token');
                });
            }, 6000)
        }).catch(() => {
        console.log("Authenticated Failed");
    })
})