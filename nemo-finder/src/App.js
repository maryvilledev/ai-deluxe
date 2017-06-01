import React, { Component } from 'react';
import {
  Button,
} from 'react-bootstrap';
import { isMobileDevice } from './util';
import axios from 'axios';

const imageSelectorId = 'image-selector';
const styles = {
  container: {
    display: 'flex',
    flexDirection: 'column',
    justifyContent: 'center',
    alignItems: 'center',
  },
  imageSelector: {
    visibility: 'hidden',
    width: '0',
    height: '0',
  },
  icon: {
    marginTop: '5px',
    maxWidth: '100%',
    maxHeight: '100%',
    borderColor: '#000',
    borderWidth: '5px',
    borderStyle: 'solid',
    borderRadius: '5px',
  }
};

class App extends Component {
  constructor(props) {
    super(props);
    this.state = { icon: '' };

    this.handleIconSelected = this.handleIconSelected.bind(this);
  }

  handleIconSelected(ev) {
    // this.setState({ icon: URL.createObjectURL(ev.target.files[0]) });

    const formData = new FormData();
    formData.append('image', ev.target.files[0])
    const config = {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    }
    axios.post('http://aideluxe.maryvilledevcenter.io/test', formData, config)
      .then(res => {
        this.setState({ icon: res.data })
      })
      .catch(err => console.log(err))

    ev.preventDefault();
  }

  displayImageSelector() {
    document.getElementById(imageSelectorId).click();
  }

  render() {
    const image = this.state.icon ?
      <img
        style={styles.icon}
        src={this.state.icon}
        alt=""
        width="500px"
      /> : null;
    return (
      <div style={styles.container}>
        <h1>Test App</h1>
        <input
          style={styles.imageSelector}
          id={imageSelectorId}
          type="file"
          accept="image/*"
          capture="camera"
          onChange={this.handleIconSelected}
        />
        <Button
          bsStyle="primary"
          onClick={this.displayImageSelector}
        >
          {isMobileDevice() ? 'Take Image' : 'Upload Image'}
        </Button>
        {image}
      </div>
    );
  }
}

export default App;
