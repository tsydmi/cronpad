<template>
  <v-card>
    <v-card-title>
      <span class="headline">Create a new Tag</span>
    </v-card-title>
    <v-card-text>
      <v-container>
        <v-form v-model="valid">
          <v-row>
            <v-text-field
                label="Name"
                v-model="newTag.name"
                @keydown.enter="saveTag"
                required
            ></v-text-field>
          </v-row>
          <v-row>
            <v-select
                label="Parent"
                hint="new tag can extend existed one"
                v-model="newTag.parent"
                item-text="name"
                :items="tags"
            ></v-select>
          </v-row>
          <v-row>
            <span class="v-label v-text-field mr-3">Color</span>
            <v-color-picker
                label="Color"
                hide-canvas
                show-swatches
                :swatches="swatches"
                swatches-max-height="400px"
                v-model="newTag.color"
            ></v-color-picker>
          </v-row>
        </v-form>
      </v-container>
    </v-card-text>
    <v-card-actions>
      <v-btn
          color="blue darken-1"
          text
          @click="saveTag"
      >
        Save
      </v-btn>
      <v-spacer></v-spacer>
      <v-btn
          color="blue darken-1"
          text
          @click="closeDialog"
      >
        Close
      </v-btn>
    </v-card-actions>
  </v-card>
</template>

<script>
import axios from 'axios'

export default {
  name: "CreateTag",
  props: {
    dialog: {
      type: Boolean,
      required: true,
    },
    tags: {
      type: Array,
      required: true,
    },
  },
  data: () => ({
    newTag: {},
    valid: true,
    swatches: [
      ['#83677b', '#445a67', '#50b4d8'],
      ['#DAAD86', '#57838d', '#96b3c2'],
      ['#659dbd', '#f3bfb3', '#ffdad1'],
      ['#5cbd95', '#ccadb2', '#b7eaf7'],
      ['#beede5', '#a7d9c9', '#6e7b8f'],
    ]
  }),
  methods: {
    saveTag() {
      let color = this.newTag.color;
      this.newTag.color = color.hex

      axios.post('/tags', this.newTag)
      .then((response) => {
        if (response.status >= 200 && response.status < 300)
        this.$emit('refreshTags', null)
        this.closeDialog()
      })
    },
    closeDialog() {
      this.$emit('closeDialog', null)
    }
  },
  watch: {
    'dialog': function () {
      this.newTag = {}
    }
  }
}
</script>

<style scoped>

</style>