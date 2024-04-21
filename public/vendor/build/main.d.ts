import './style/main.scss';
import * as bootstrap from 'bootstrap';
import Swal from 'sweetalert2';
import ckeditor from './ckeditor';
import htmx from './htmx';
declare global {
    interface Window {
        htmx: typeof htmx;
        ckeditor: typeof ckeditor;
        bootstrap: typeof bootstrap;
        Swal: typeof Swal;
    }
}