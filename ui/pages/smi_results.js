import React, { useEffect } from 'react';
import MeshplaySMIResults from '../components/MeshplaySMIResults';
import { updatepagepath } from '../lib/store';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import Head from 'next/head';
import { getPath } from '../lib/path';
import { Paper, withStyles } from '@material-ui/core';

const styles = { paper: { maxWidth: '90%', margin: 'auto', overflow: 'hidden' } };

/**
 * @deprecated This functionality has been deprecated and its child components can be left behind
 */
const SMIResults = (props) => {
  useEffect(() => {
    props.updatepagepath({ path: getPath() });
  }, [props.updatepagepath]);

  return (
    <React.Fragment>
      <Head>
        <title>SMI Results | Meshplay</title>
      </Head>
      <Paper className={props.classes.paper}>
        <MeshplaySMIResults />
      </Paper>
    </React.Fragment>
  );
};

const mapDispatchToProps = (dispatch) => ({
  updatepagepath: bindActionCreators(updatepagepath, dispatch),
});

export default withStyles(styles)(connect(null, mapDispatchToProps)(SMIResults));
