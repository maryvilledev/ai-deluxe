import React, { Component } from 'react';
import {
  Button,
} from 'react-bootstrap';
import {
  isMobileDevice,
  capitalize,
} from './util';
import axios from 'axios';
import Loader from 'halogen/BounceLoader';

const API_URL = process.env.REACT_APP_API_URL;

const getCharacter = () => {
  const { pathname } = window.location;
  const character = pathname !== '/' ? pathname.substring(1) : 'Nemo';
  return capitalize(character)
}

const imageSelectorId = 'image-selector';
const styles = {
  container: {
    height: '100vh',
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'center',
  },
  p: {
    maxWidth: '500px',
    padding: '10px',
  },
  loader: {
    height: '50%',
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
    margin: '5px',
    marginTop: '20px',
    borderColor: '#999999',
    borderWidth: '5px',
    borderStyle: 'solid',
    borderRadius: '5px',
  },
  img: {
    maxWidth: '100%',
    maxHeight: '100%',
  }
};

class App extends Component {
  constructor(props) {
    super(props);
    this.state = {
      icon: '' ,
      isUploading: false,
    };

    this.handleIconSelected = this.handleIconSelected.bind(this);
  }

  handleIconSelected(ev) {
    this.setState({ isUploading: true });
    const formData = new FormData();
    formData.append('image', ev.target.files[0])
    const config = {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    }
    axios.post(`${API_URL}/find/${getCharacter()}`, formData, config)
      .then(res => {
        this.setState({
          icon: res.data,
          isUploading: false,
        });
      })
      .catch(err => {
        this.setState({ isUploading: false });
        alert('An error occurred while uploading image.')
        console.log(err)
      })

    ev.preventDefault();
  }

  displayImageSelector() {
    document.getElementById(imageSelectorId).click();
  }

  render() {
    const {
      icon,
      isUploading,
    } = this.state;
    const image = icon && !isUploading ?
      <div style={styles.icon}>
        <img
          style={styles.img}
          src={this.state.icon}
          alt=""
          width="500px"
        />
      </div> : null;
    const loader = isUploading ?
      <div style={styles.loader}>
        <h3>Loading...</h3>
        <Loader
          color="#337AB7"
          size="50px"
          margin="4px"
        />
      </div> : null;

    return (
      <div style={styles.container}>
        <h1>Find {getCharacter()}!</h1>
        <p style={styles.p}>
          Upload an image and our advanced, sentient AI will locate and outline {getCharacter()}. If the image doesn't contain {getCharacter()}, the AI will drop an "X" on the image.
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
          disabled={isUploading}
          onClick={this.displayImageSelector}
        >
          {isMobileDevice() ? 'Take Image' : 'Upload Image'}
        </Button>
        {loader}
        {image}
      </div>
    );
  }
}

export default App;
