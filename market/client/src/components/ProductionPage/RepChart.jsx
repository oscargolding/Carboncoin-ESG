import React from 'react';
import PropTypes from 'prop-types';
import Chart from 'react-google-charts';

/**
 * Breakdown of data on the blockchain.
 * @returns a Chart representing the breakdown on the blockchain.
 */
const RepChart = (props) => {
  const { environment, social, governance, } = props;
  console.log(`${environment} ${social} ${governance}`);
  return (
    <Chart
      width={'500px'}
      height={'400px'}
      chartType="BarChart"
      loader={< div > Loading Chart</div>}
      data={
        [
          [
            'Element',
            'Reputation',
            { role: 'style', },
            {
              sourceColumn: 0,
              role: 'annotation',
              type: 'string',
              calc: 'stringify',
            }
          ],
          ['Environment', environment, 'green', null],
          ['Social', social, 'pink', null],
          ['Governance', governance, 'grey', null]
        ]}
      options={{
        title: 'Reputation Categories for Producer',
        width: 600,
        height: 400,
        bar: { groupWidth: '95%', },
        legend: { position: 'none', },
      }}
      // For tests
      rootProps={{ 'data-testid': '6', }}
    />
  );
};

export default RepChart;

RepChart.propTypes = {
  environment: PropTypes.number.isRequired,
  social: PropTypes.number.isRequired,
  governance: PropTypes.number.isRequired,
};
