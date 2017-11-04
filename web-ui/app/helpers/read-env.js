import { helper } from '@ember/component/helper';
import ENV from '../config/environment';

export function readEnv(params/*, hash*/) {
  if (params.length === 0) {
    return undefined;
  }

  return Ember.Object.create(ENV).get(params[0]);
}

export default helper(readEnv);
