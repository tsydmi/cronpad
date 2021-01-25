<template>
  <v-dialog
      :value="value"
      @input="emitChange"
      max-width="400px"
  >
    <v-card>
      <v-card-title>
        <span class="headline">Create a new Tag</span>
      </v-card-title>
      <v-card-text>
        <v-container>
          <v-form ref="form" v-model="valid">
            <v-row>
              <v-text-field
                  label="Name"
                  v-model="newTag.name"
                  @keydown.enter="saveTag"
                  :rules="rules.name"
              ></v-text-field>
            </v-row>
            <v-row>
              <v-text-field
                  label="Description"
                  v-model="newTag.description"
                  @keydown.enter="saveTag"
              ></v-text-field>
            </v-row>
            <v-row>
              <tag-color-picker v-model="newTag.color"/>
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
            @click="emitChange(false)"
        >
          Close
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script>
import TagService from "@/service/TagService"
import TagColorPicker from "@/components/tags/TagColorPicker"

export default {
  name: "CreateTagDialog",
  components: {
    TagColorPicker,
  },

  props: {
    value: {
      type: Boolean,
      required: false,
    },
  },
  data: () => ({
    valid: true,
    newTag: {},
    rules: {
      name: [
        v => !!v || 'Name is required',
      ],
    },
  }),
  methods: {
    saveTag() {
      if (this.$refs.form.validate()) {
        TagService.create(this.newTag)
            .then(() => {
              this.$emit('refreshTags', null)
              this.emitChange(false)
            })
      }
    },
    emitChange(value) {
      this.$emit('input', value)
    }
  },
  watch: {
    'value': function () {
      if (this.value === false) {
        this.$refs.form.reset()
      }
    },
  },
}
</script>