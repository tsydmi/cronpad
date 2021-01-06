import Vue from 'vue'
import App from './App.vue'
import vuetify from './plugins/vuetify';
import axios from 'axios'
import Keycloak from 'keycloak-js'
import './sass/variables.scss'

import router from "@/config/router";
import {initAxiosInterceptors} from "@/config/axios"

Vue.config.productionTip = false

Vue.prototype.$http = axios
axios.defaults.baseURL = process.env.VUE_APP_BACKEND_URL

let initOptions = {
    url: process.env.VUE_APP_KEYCLOAK_URL, realm: 'cronpad', clientId: 'vue-frontend', onLoad: 'login-required'
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
            keycloak.updateToken(70).then((refreshed) => {
                if (refreshed) {
                    console.log('Token refreshed' + refreshed);
                } else {
                    var expirationTime = Math.round(keycloak.tokenParsed.exp + keycloak.timeSkew - new Date().getTime() / 1000);
                    console.log(`Token not refreshed, valid for ${expirationTime} seconds`);
                }
            }).catch(() => {
                console.log('Failed to refresh token');
            });
        }, 6000)
    }).catch(() => {
    console.log("Authenticated Failed");
})