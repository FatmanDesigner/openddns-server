import DS from 'ember-data';

export default DS.RESTAdapter.extend({
  namespace: '/api/rest',
  headers: {
    Authorization: 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJnaXRodWI6MzkyNjM0MCJ9.CPJ9LU6hO4hWvLH8tMbu70qEGSIX-OOoOWJkSID8Ao0'
  }
});
