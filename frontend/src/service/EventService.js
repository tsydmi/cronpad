import axios from "axios";

export default {
    create(event) {
        let body = {
            start: new Date(event.start).toISOString(),
            end: new Date(event.end).toISOString(),
            name: event.name,
            tag: event.tag,
            timed: event.timed
        }

        return axios.post('/events', body)
    },
    update(event) {
        let body = {
            start: new Date(event.start).toISOString(),
            end: new Date(event.end).toISOString(),
            name: event.name,
            tag: event.tag.id,
            timed: event.timed
        }

        return axios.put(`/events/${event.id}`, body)
    },
    delete(event) {
        return axios.delete(`/events/${event.id}`)
    },
}