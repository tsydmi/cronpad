<template>
  <v-dialog
      :value="value"
      @input="emitChange"
      max-width="600px"
  >
    <v-card>
      <v-card-title>
        <span class="headline">Edit Project</span>
      </v-card-title>
      <v-card-text>
        <v-container>
          <v-form
              ref="form"
              lazy-validation
          >
            <v-row>
              <v-text-field
                  label="Name"
                  v-model="project.name"
                  :rules="rules.name"
                  required
              ></v-text-field>
            </v-row>

            <v-row>
              <v-text-field
                  label="Description"
                  v-model="project.description"
                  :rules="rules.description"
                  required
              ></v-text-field>
            </v-row>

            <v-row>
              <user-select v-model="project.users" :users="users"/>
            </v-row>

            <v-row>
              <date-picker label="From" v-model="project.start"/>
            </v-row>

            <v-row>
              <date-picker label="To" v-model="project.end"/>
            </v-row>
          </v-form>
        </v-container>
      </v-card-text>

      <div class="d-flex justify-center">
        <div class="pl-3 pr-3 global-form-error error--text">
          {{ globalFormError }}
        </div>
      </div>

      <v-card-actions>
        <v-btn
            color="primary"
            text
            @click="update"
        >
          Update
        </v-btn>
        <v-spacer></v-spacer>
        <v-btn
            color="primary"
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
import ProjectService from "@/service/ProjectService"
import UserSelect from "@/components/UserSelect"
import DatePicker from "@/components/DatePicker"

export default {
  name: "EditProjectDialog",

  components: {
    UserSelect,
    DatePicker,
  },

  props: {
    value: {
      type: Boolean,
      required: false,
    },
    project: {
      type: Object,
      required: false,
    },
    users: {
      type: Array,
      required: true,
    },
  },
  data: () => ({
    globalFormError: '',
    rules: {
      name: [
        v => !!v || 'Name is required',
      ],
      description: [
        v => !!v || 'Description is required',
      ],
    },
  }),
  methods: {
    update() {
      if (this.$refs.form.validate()) {
        ProjectService.update(this.project)
            .then(() => {
              this.$emit('refreshProjects', null)
              this.emitChange(false)
            })
            .catch(error => {
              if (error && error.response && (error.response.status === 400 || error.response.status === 404)) {
                this.globalFormError = error.response.data.error
              }
            })
      }
    },
    emitChange(value) {
      this.$emit('input', value)
    },
  },
}
</script>