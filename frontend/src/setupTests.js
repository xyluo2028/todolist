import { TextEncoder } from 'util';

global.TextEncoder = TextEncoder;
global.IS_REACT_ACT_ENVIRONMENT = true;
