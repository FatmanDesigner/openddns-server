import Route from '@ember/routing/route';
import Ember from 'ember';

export default Route.extend({
  model() {
    return Ember.$.ajax({
      url: `/api/rest/appInfo`,
      dataType: 'json',
      beforeSend: function(xhr){
        if (!document.cookie.length) {
          return;
        }

        var rgx = /accessToken=(.+);?/g;
        var match = rgx.exec(document.cookie);
        if (match && match[1]) {
          xhr.setRequestHeader('Authorization', 'Bearer ' + match[1]);
        } else {
          return;
        }
      }
    }).promise();
  }
});
