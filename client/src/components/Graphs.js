import { Line, Bar } from 'react-chartjs-2';
import {
  Chart as ChartJS,
  BarElement,
  CategoryScale,
  Legend,
  LineElement,
  LinearScale,
  PointElement,
  TimeScale,
  Tooltip
} from 'chart.js';
import zoomPlugin from 'chartjs-plugin-zoom';
import 'chartjs-adapter-date-fns';

ChartJS.register(
    BarElement,
    CategoryScale,
    Legend,
    LineElement,
    LinearScale,
    PointElement,
    TimeScale,
    Tooltip,
    zoomPlugin
);

export function graphColors() {
  const rootStyles = getComputedStyle(document.documentElement);
  const graphColorA = rootStyles.getPropertyValue('--color-graphA').trim();
  const graphColorB = rootStyles.getPropertyValue('--color-graphB').trim();
  const graphColorC = rootStyles.getPropertyValue('--color-graphC').trim();

  return { graphColorA, graphColorB, graphColorC }
}

export function getDefaultChartOptions(additionalOptions = {}, data = null) {
  const allDataPoints = data && data.datasets
    ? data.datasets.flatMap(dataset => dataset.data)
    : [];
  const maxYValue = allDataPoints.length > 0
    ? Math.max(...allDataPoints) + 1
    : undefined;

  return {
    responsive: true,
    maintainAspectRatio: false,
    plugins: {
      legend: {
        display: true,
        position: 'right',
      },
      zoom: {
        zoom: {
          wheel: { enabled: true },
          pinch: { enabled: true },
          mode: 'x',
        },
        pan: {
          enabled: true,
          mode: 'x',
        },
      },
    },
    scales: {
      y: {
        beginAtZero: true,
        suggestedMax: maxYValue,
      },
    },
    ...additionalOptions,
  };
}

export function LineChart({ data, options, style }) {
  const defaultOptions = getDefaultChartOptions(options, data);
  return <Line data={data} options={defaultOptions} style={style} />;
}

export function BarChart({ data, options, style }) {
  const defaultOptions = getDefaultChartOptions(options, data);
  return <Bar data={data} options={defaultOptions} style={style} />;
}
