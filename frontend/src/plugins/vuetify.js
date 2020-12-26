import Vue from 'vue';
import Vuetify from 'vuetify/lib';

Vue.use(Vuetify);

export default new Vuetify({
    theme: {
        options: { customProperties: true },
        themes: {
            light: {
                primary: '#32c8c8',
                secondary: '#FFFF00',
                contrast: '#123123',
                inactive: '#b0bec5',
                error: '#b71c1c',
            },
        }
    },
});
