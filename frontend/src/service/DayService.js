import axios from "axios";

export default {
    findByDayRange(from, to) {
        let firstDay = from.toISOString().split('T')[0]
        let lastDay = to.toISOString().split('T')[0]
        return axios.get(`/days?from=${firstDay}&to=${lastDay}`)
    },
}