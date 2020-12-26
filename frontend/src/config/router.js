import Vue from 'vue'
import VueRouter from 'vue-router'

import Home from "@/pages/Home";
import WorkingTime from "@/pages/WorkingTime";
import Settings from "@/pages/Settings";
import About from "@/pages/About";
import Page404 from "@/pages/Page404";
import Projects from "@/pages/Projects";

const router = new VueRouter({
    mode: 'history',
    routes: [
        {path: '/', component: Home},
        {path: '/working-time', component: WorkingTime},
        {path: '/projects', component: Projects},
        {path: '/settings', component: Settings},
        {path: '/about', component: About},
        {path: "*", component: Page404}
    ]
});

Vue.use(VueRouter)

export default router