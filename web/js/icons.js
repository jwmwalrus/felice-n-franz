import { library, dom, icon } from '@fortawesome/fontawesome-svg-core';
import {
    faBroom,
    faCheck,
    faCompress,
    faCopy,
    faCrop,
    faEllipsisH,
    faEye,
    faMinus,
    faMinusCircle,
    faPause,
    faPlay,
    faPlus,
    faQuestion,
    faRss,
    faSearch,
    faShoppingBag,
    faStop,
    faSync,
    faTimesCircle,
} from '@fortawesome/free-solid-svg-icons';

library.add(
    faBroom,
    faCheck,
    faCompress,
    faCopy,
    faCrop,
    faEllipsisH,
    faEye,
    faMinus,
    faMinusCircle,
    faPause,
    faPlay,
    faPlus,
    faQuestion,
    faRss,
    faSearch,
    faShoppingBag,
    faStop,
    faSync,
    faTimesCircle,
);

dom.watch();

const broomIcon = icon({ prefix: 'fas', iconName: 'broom' }).html;
const compressIcon = icon({ prefix: 'fas', iconName: 'compress' }).html;
const copyIcon = icon({ prefix: 'fas', iconName: 'copy' }).html;
const ellipsisHIcon = icon({ prefix: 'fas', iconName: 'ellipsis-h' }).html;
const pauseIcon = icon({ prefix: 'fas', iconName: 'pause' }).html;
const playIcon = icon({ prefix: 'fas', iconName: 'play' }).html;
const shoppingBagIcon = icon({ prefix: 'fas', iconName: 'shopping-bag' }).html;

export {
    broomIcon,
    compressIcon,
    copyIcon,
    ellipsisHIcon,
    pauseIcon,
    playIcon,
    shoppingBagIcon,
};
