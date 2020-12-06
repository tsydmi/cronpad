import Vue from 'vue'
import VueRouter from 'vue-router'

import Home from "@/pages/Home";
import WorkingTime from "@/pages/WorkingTime";
import Settings from "@/pages/Settings";
import About from "@/pages/About";
import Page404 from "@/pages/Page404";

const router = new VueRouter({
    mode: 'history',
    routes: [
        {path: '/', component: Home},
        {path: '/working-time', component: WorkingTime},
        {path: '/settings', component: Settings},
        {path: '/about', component: About},
        {path: "*", component: Page404}
]
});

Vue.use(VueRouter)

export default router