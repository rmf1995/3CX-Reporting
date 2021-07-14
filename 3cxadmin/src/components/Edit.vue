<template>
   <div>
      <v-card>
         <v-card-title>
            Edit 3CX Server
         </v-card-title>
         <form>
            <v-text-field
               v-model="item.Name"
               :counter="254"
               label="Name"
               required
            ></v-text-field>
            <v-text-field
               v-model="item.Location"
               label="Location"
               required
            ></v-text-field>
            <v-text-field
               v-model="item.URL"
               label="URL"
               required
            ></v-text-field>
            <v-text-field
               v-model="item.UserName"
               label="UserName"
               required
            ></v-text-field>
            <v-text-field
               v-model="item.PwStateID"
               label="PwState ID"
               required
            ></v-text-field>
            <v-text-field
               v-model="item.CustomerID"
               label="Customer ID"
               required
            ></v-text-field>
            <v-text-field
               v-model="item.Bespoke"
               label="Bespoke"
            ></v-text-field>
            <v-btn
               class="mr-4"
               @click="updateItem"
               >
               submit
            </v-btn>
         </form>
      </v-card>
   </div>
</template>




<script>
export default{
        data(){
            return{
                item:{}
            }
        },

        created: function(){
            this.getItem();
        },

        methods: {
            getItem()
            {
              let uri = 'https://URL/v1/3cx/edit/' + this.$route.params.id;
                this.axios.get(uri).then((response) => {
                    this.item = response.data;
                });
            },

            updateItem()
            {
              let uri = 'https://URL/v1/3cx/update/' + this.$route.params.id;
                this.axios.post(uri, this.item).then((response) => {
                  this.$router.push({name: 'Index'});
                });
            }
        }
    }
</script>