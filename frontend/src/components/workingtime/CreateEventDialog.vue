<template>
  <v-dialog
      :value="value"
      @input="emitChange"
      max-width="600px"
  >
    <v-card>
      <v-card-title>
        <span class="headline">Create a new Event</span>
      </v-card-title>

      <v-form
          ref="form"
          v-model="valid"
          lazy-validation
          @submit.prevent="saveEvent"
      >
        <v-card-text>
          <v-text-field
              autofocus
              label="NAME"
              v-model="name"
              :rules="rules.name"
              required
          ></v-text-field>
        </v-card-text>

        <v-card-actions>
          <v-btn
              class="ma-2"
              color="primary"
              plain
              type="submit"
          >
            Save
          </v-btn>
          <v-btn
              class="ma-2"
              color="inactive--text"
              text
              @click="emitChange(false)"
          >
            Cancel
          </v-btn>
        </v-card-actions>
      </v-form>
    </v-card>
  </v-dialog>
</template>

<script>
import EventService from '@/service/EventService'

export default {
  name: "DeleteDialog",

  props: {
    value: {
      type: Boolean,
      required: false,
    },
    event: {
      type: Object,
      required: false,
    },
  },
  data: () => ({
    name: "",
    valid: true,
    rules: {
      name: [
        v => !!v || 'Name is required',
      ],
    },
  }),
  methods: {
    saveEvent() {
      if (this.$refs.form.validate()) {

        const event = {
          name: this.name,
          start: this.event.start,
          end: this.event.end,
          tag: this.event.tag,
          project: this.event.project,
          timed: this.event.timed,
        }

        EventService.create(event)
            .then(() => {
              this.emitChange(false)
            })
      }
    },
    emitChange(value) {
      this.$emit('refreshEvents', null)
      this.$emit('input', value)
    },
  },
  watch: {
    'value': function () {
      this.$refs.form.reset()
    }
  },
}
</script>