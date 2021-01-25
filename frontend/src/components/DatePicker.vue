<template>
  <v-row>
    <v-col cols="2" class="d-flex align-center">
      <span class="date-picker-label">{{ label }}</span>
    </v-col>
    <v-col cols="9">
      <v-menu
          v-model="menu"
          :close-on-content-click="false"
          transition="scale-transition"
          offset-y
          min-width="auto"
      >
        <template v-slot:activator="{ on, attrs }">
          <v-text-field
              :value="dateToText"
              @input="emitChange"
              @change="$emit('change', '')"
              prepend-icon="mdi-calendar"
              v-bind="attrs"
              @blur="date = parseDate(value)"
              v-on="on"
          ></v-text-field>
        </template>
        <v-date-picker
            :value="dateToText"
            @input="emitChange"
            @change="$emit('change', '')"
            offset-y
            min-width="auto"
            no-title
        ></v-date-picker>
      </v-menu>
    </v-col>
    <v-col cols="1" class="d-flex align-center">
      <v-btn
          icon
          @click="reset"
      >
        <v-icon>far fa-times-circle</v-icon>
      </v-btn>
    </v-col>
  </v-row>
</template>

<script>
import dayjs from "dayjs"

export default {
  name: "DatePicker",

  props: {
    value: {
      type: String,
      required: false,
    },
    label: {
      type: String,
      required: false,
    },
  },
  data: () => ({
    menu: false,
  }),
  methods: {
    parseDate(date) {
      if (!date) {
        return null
      }

      return dayjs(date, 'M/D/YYYY').hour(0).minute(0).second(0)
    },
    emitChange(value) {
      this.menu = false

      if (!value) {
        this.$emit('input', value)
      } else {
        this.$emit('input', dayjs(value, 'M/D/YYYY').hour(0).minute(0).second(0).toISOString())
      }
    },
    reset() {
      this.$emit('input', null)
    },
  },
  computed: {
    dateToText() {
      if (!this.value) {
        return null
      }
      return dayjs(this.value).format('YYYY-MM-DD');
    },
  },
}
</script>

<style>
.date-picker-label {
  font-size: 16px;
}
</style>