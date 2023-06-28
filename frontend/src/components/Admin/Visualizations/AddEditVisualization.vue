<!--
<template>
  <div class="modal-card" style="width:800px">
    <section class="info-tiles" style="margin-top:50px">
      <article class="tile is-child box">
        <h1 style="font-size:50px;margin-bottom:20px; text-align:center">Visualization</h1>
        <b-steps
        style="min-height:400px;width:100%;color:'red'"
        v-model="activeStep"
        :has-navigation="false">
          <b-step-item style="height:300px" label="Type">
            <div class="has-text-centered">
              <h1 class="title has-text-centered">Select Display Type</h1>
              <p>{{ view.class == undefined ? "none" : view.class }} selected</p>
              <b-button
                class="type-button"
                size="is-large"
                icon-left="chart-line"
                @click="setViewClass('line')"
                :type="view.class == 'line' ? 'is-primary' : 'none'"
              >Line</b-button>
              <b-button
                class="type-button"
                size="is-large"
                icon-left="chart-bar"
                style="margin-left:50px"
                @click="setViewClass('bar')"
                :type="view.class == 'bar' ? 'is-primary' : 'none'"
              >Bar</b-button>
              <br />
              <b-button
                class="type-button"
                size="is-large"
                icon-left="chart-pie"
                @click="setViewClass('pie')"
                :type="view.class == 'pie' ? 'is-primary' : 'none'"
              >Pie</b-button>
              <b-button
                class="type-button"
                size="is-large"
                icon-left="table"
                style="margin-left:50px"
                @click="setViewClass('table')"
                :type="view.class == 'table' ? 'is-primary' : 'none'"
              >Table</b-button>
            </div>
          </b-step-item>

          <b-step-item style="height:300px" label="Information">
              <Information ref="info" :viewClass="view.class" :visualizationToEdit="visualizationToEdit"></Information>
          </b-step-item>
          <b-step-item style="height:300px" label="Finish" disabled>
            <h1 class="title has-text-centered">Finish</h1>
            <div class="columns is-mobile is-centered">
              <div class="column is-half">
                <b-field label="Visualization Name">
                  <b-input v-model="view.name" placeholder="Enter a name..." required></b-input>
                </b-field>
                <b-button
                  @click="saveView"
                  type="is-primary"
                  style="float:right;width: 100%"
                ><b-icon icon="check"></b-icon><span>Submit</span></b-button>
              </div>
            </div>
            </b-step-item>
            <div class="columns">
              <div class="column">
                <b-button
                  outlined
                  type="is-danger"
                  icon-pack="fas"
                  icon-left="backward"
                  @click="stepBackward"
                  style="float:left"
                  v-if="activeStep > 0">
                  Previous
              </b-button>
              <b-button
                  outlined
                  type="is-success"
                  icon-pack="fas"
                  icon-right="forward"
                  @click="stepForward"
                  style="float: right"
                  v-if="activeStep < 2">
                  Next
              </b-button>
            </div>
          </div>
        </b-steps>
      </article>
    </section>
  </div>
</template>
-->

<template>
  <div class="modal-card" style="width:800px">
    <section class="info-tiles" style="margin-top: 50px">
      <article class="tile is-child box">
        <h1 style="font-size: 25px; margin-bottom: 20px; text-align: center">Visualization</h1>
        <b-steps
          style="min-height: 200px; width: 100%; color: 'red'"
          :has-navigation="true"
        >
          <b-step-item style="min-height: 100px" label="Type">
            <div class="has-text-centered">
              <h1 class="title has-text-centered" style="font-size: 25px">Select Display Type</h1>
              <div class="type-buttons">
                <b-button
                  class="type-button"
                  size="is-medium"
                  icon-left="chart-line"
                  @click="setViewClass('line')"
                  :type="view.class == 'line' ? 'is-primary' : 'none'"
                >Line</b-button>
                <b-button
                  class="type-button"
                  size="is-medium"
                  icon-left="chart-bar"
                  @click="setViewClass('bar')"
                  :type="view.class == 'bar' ? 'is-primary' : 'none'"
                >Bar</b-button>
                <b-button
                  class="type-button"
                  size="is-medium"
                  icon-left="chart-pie"
                  @click="setViewClass('pie')"
                  :type="view.class == 'pie' ? 'is-primary' : 'none'"
                >Pie</b-button>
                <b-button
                  class="type-button"
                  size="is-medium"
                  icon-left="table"
                  @click="setViewClass('table')"
                  :type="view.class == 'table' ? 'is-primary' : 'none'"
                >Table</b-button>
              </div>
            </div>
          </b-step-item>

          <b-step-item style="min-height: 100px" label="Information">
            <Information
              ref="info"
              :viewClass="view.class"
              :visualizationToEdit="visualizationToEdit"
            ></Information>
          </b-step-item>
          <b-step-item style="min-height: 100px" label="Finish" disabled>
            <h1 class="title has-text-centered">Finish</h1>
            <div class="columns is-mobile is-centered">
              <div class="column is-half">
                <b-field label="Visualization Name">
                  <b-input v-model="view.name" placeholder="Enter a name..." required></b-input>
                </b-field>
                <b-button
                  @click="saveView"
                  type="is-primary"
                  style="float: right; width: 100%"
                ><b-icon icon="check"></b-icon><span>Submit</span></b-button>
              </div>
            </div>
          </b-step-item>
        </b-steps>
      </article>
    </section>
  </div>
</template>

<style scoped>
@media (max-width: 768px) {
  .modal-card {
    width: 90%;
  }
}

.type-buttons {
  display: flex;
  flex-wrap: wrap;
  justify-content: center;
  margin-top: 20px;
}

.type-button {
  margin: 5px;
}

.navigation-buttons {
  margin-top: 20px;
  display: flex;
  justify-content: space-between;
}

@media (max-width: 768px) {
  .navigation-buttons {
    flex-direction: column;
  }
}
</style>

<script>
import Information from './views/Information'
export default {
  props: ['visualizationToEdit'],
  data() {
    return {
      view: {
        name: "",
        class: "",
        selectedIndex: {},
        index: "",
        fields: [],
        fieldNames: []
      },
      activeStep: 0,
      allFields: []
    };
  },
  mounted() {
    // if we passed an object to edit
    if (this.visualizationToEdit.class) {
      this.view.name = this.visualizationToEdit.name
      this.view.class = this.visualizationToEdit.class
      this.view.index = this.visualizationToEdit.index
      this.view.fields = this.visualizationToEdit.fields
      this.view.fieldNames = this.visualizationToEdit.fieldNames
    }
  },
  methods: {
    changedView(view) {
      this.view = view
    },
    setViewClass(type) {
      this.view.class = type
    },
    saveView() {
      switch (this.view.class) {
        case 'line':
          this.view.fields = this.view.fields.slice(0, 2)
          this.view.fieldNames = this.view.fieldNames.slice(0, 2)
          break;
        case 'bar':
        case 'pie':
          this.view.fields = this.view.fields.slice(0, 1)
          this.view.fieldNames = this.view.fieldNames.slice(0, 1)
          break;
        default:
          break;
      }
      if (this.visualizationToEdit.class) {
        const body = {
          uuid: this.visualizationToEdit.uuid,
          name: this.view.name,
          class: this.view.class,
          index: this.view.selectedIndex.index,
          fields: this.view.fields,
          fieldNames: this.view.fieldNames
        }
        fetch('/api/view/update', {
          method: 'post',
          body: JSON.stringify(body)
        })
          .then(response => {
            if (response.ok) {
              return Promise.all([response.ok, response.json()]);
            } else {
              return Promise.all([response.ok, response.text()]);
            }
          })
          .then(data => {
            if (data[0] === false) {
              this.$buefy.snackbar.open(data[1]);
            } else {
              this.visualizationToEdit.name = this.view.name
              this.visualizationToEdit.class = this.view.class
              this.visualizationToEdit.index = this.view.index
              this.visualizationToEdit.fields = this.view.fields
              this.visualizationToEdit.fieldNames = this.view.fieldNames
              this.$parent.close()
              this.$emit('editedVisualization', this.view)
            }
          });
      } else {
        const body = {
          name: this.view.name,
          class: this.view.class,
          index: this.view.selectedIndex.index,
          fields: this.view.fields,
          fieldNames: this.view.fieldNames
        }
        fetch('/api/view/add', {
          method: 'post',
          body: JSON.stringify(body)
        })
          .then(response => {
            if (response.ok) {
              return Promise.all([response.ok, response.json()]);
            } else {
              return Promise.all([response.ok, response.text()]);
            }
          })
          .then(data => {
            if (data[0] === false) {
              this.$buefy.snackbar.open(data[1]);
            } else {
              this.$parent.close()
              this.$emit('addedVisualization', this.view)
            }
          });
      }
    },
    stepForward() {
      if (this.activeStep < 2) {
        this.activeStep++
      }
      if (this.activeStep === 2) {
        this.view = this.$refs.info.getView()
      }
    },
    stepBackward() {
      if (this.activeStep !== 0) {
        this.activeStep--
      }
    }
  },
  components: {
    Information
  }
};
</script>

<style scoped lang="scss">
.type-button {
  width: 150px;
}
.button {
  margin-top: 20px;
  margin-bottom: 20px;
}
</style>
