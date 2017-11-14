import Controller from '@ember/controller';
import Ember from 'ember';

export default Controller.extend({
  generateSecretAndHide(appid) {
    Ember.$.ajax({
      url: `/api/generate-secret?appid=${appid}`,
      dataType: 'json',
      beforeSend
    }).promise()
    .then(({secret}) => {
      this.set('secret', secret);

      setTimeout(() => {
        this.set('secret', null);
      }, 5000);
    });
  }
});

function beforeSend(xhr){
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
