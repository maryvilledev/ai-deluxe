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
  p: {
    maxWidth: '500px',
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
    const formData = new FormData();
    formData.append('image', ev.target.files[0])
    const config = {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    }
    axios.post('http://aideluxe.maryvilledevcenter.io:8080/test', formData, config)
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
        <h1>Find Nemo!</h1>
        <p style={styles.p}>
          Upload an image and our advanced, sentient AI will locate and outline Nemo. If the image doesn't contain Nemo, the AI will drop an "X" on the image.
        </p>
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
