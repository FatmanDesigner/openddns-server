import Ember from 'ember';
import DS from 'ember-data';

export default DS.RESTAdapter.extend({
  namespace: '/api/rest',
  headers: Ember.computed(function () {
    if (!document.cookie.length) {
      return '';
    }

    var rgx = /accessToken=(.+);?/g;
    var match = rgx.exec(document.cookie);
    if (match && match[1]) {
      return {
        Authorization: 'Bearer ' + match[1]
      }
    } else {
      return {};
    }
  })
});
