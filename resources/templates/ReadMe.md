## App ressources

### Yaml 
The yaml templates are used to deploy a service instance in the cluster. 
Remember to change the provider dtos aswell if you change the yaml template.

### Schema forms
NGX schema forms is a json schema which is capable of generating forms.
The backend has to deliver such a json form to tell the frontend which data is required to deploy a instance of a service.
Every provider must be able to return a json form. 
In this template folder you can find creation forms for each provider called `create_form.json`.
Since some form values are set during runtime i.e. the name of a service. Remember to change the provider dtos aswell if you change a creation form.