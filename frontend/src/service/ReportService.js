import axios from "axios";

export default {
    findAll(query) {
        return axios.post('/admin/reports', query)
    },
}