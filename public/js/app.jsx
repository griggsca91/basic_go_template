import React from 'react';
import { render } from 'react-dom';
import Button from 'material-ui/Button';

function App() {
    return (
        <Button raised color="primary">
            Hello World
        </Button>
    );
}

render(<App />, document.querySelector('#app'));
