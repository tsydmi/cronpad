import axios from "axios"
import cloneDeep from 'clone-deep'

export default {
    get() {
        return axios.get('/settings')
    },
    update(settings) {
        return axios.put('/settings', this.convertToBody(settings))
    },

    convertToBody(s) {
        let body = cloneDeep(s)

        body.timeRange.min = parseInt(body.timeRange.min)
        body.timeRange.max = parseInt(body.timeRange.max)
        return body
    },
}