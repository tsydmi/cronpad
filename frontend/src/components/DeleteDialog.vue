<template>
  <v-dialog
      max-width="600px"
      :value="value"
      @input="emitChange"
  >
    <v-sheet
        class="px-7 pt-7 pb-4 mx-auto text-center d-inline-block"
        color="blue-grey darken-3"
        dark
    >
      <div class="grey--text text--lighten-1 text-body-2 mb-4">
        Are you sure you want to delete {{ itemName }}?
      </div>
      <v-btn
          class="ma-2"
          color="error"
          plain
          @click="remove"
      >
        Delete
      </v-btn>
      <v-btn
          class="ma-2"
          color="inactive--text"
          text
          @click="emitChange(false)"
      >
        Cancel
      </v-btn>
    </v-sheet>
  </v-dialog>
</template>

<script>
export default {
  name: "DeleteDialog",

  props: {
    value: {
      type: Boolean,
      required: false,
    },
    itemName: {
      type: String,
      required: false,
    },
  },
  methods: {
    remove() {
      this.$emit('delete', '')
    },
    emitChange(value) {
      this.$emit('input', value)
    },
  },
}
</script>