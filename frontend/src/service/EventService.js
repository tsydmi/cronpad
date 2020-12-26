import axios from "axios";

export default {
    create(event) {
        return axios.post('/events', this.convertEventToRequestBody(event))
    },
    update(event) {
        return axios.put(`/events/${event.id}`, this.convertEventToRequestBody(event))
    },
    delete(event) {
        return axios.delete(`/events/${event.id}`)
    },
    convertEventToRequestBody(event) {
        return {
            start: new Date(event.start).toISOString(),
            end: new Date(event.end).toISOString(),
            name: event.name,
            tag: event.tag ? event.tag.id : null,
            project: event.project ? event.project.id : null,
            timed: event.timed
        }
    },
}