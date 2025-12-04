import React from 'react';
import { createRoot } from 'react-dom/client';
import { DisplayComp } from '../../page/control';
import '@douyinfe/semi-ui/dist/css/semi.min.css';
import '../../styles/semi-theme-overrides.scss';

const container = document.createElement('div');
container.id = 'app';
container.style.minHeight = '150px';
document.body.appendChild(container);
const root = createRoot(container);
root.render(<DisplayComp />);
