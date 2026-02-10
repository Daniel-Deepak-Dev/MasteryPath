import React from 'react';
import { View, Text, StyleSheet } from 'react-native';
import Svg, { Polygon, Line, Circle, Text as SvgText } from 'react-native-svg';
import { RadarAxisColors } from '@/constants/theme';

interface RadarDataPoint {
  label: string;
  value: number; // 0-100
}

interface RadarChartProps {
  data: RadarDataPoint[];
  size?: number;
  fillColor?: string;
  strokeColor?: string;
}

const LEVELS = 4; // Number of concentric rings

export default function RadarChart({
  data,
  size = 200,
  fillColor = 'rgba(108, 99, 255, 0.25)',
  strokeColor = '#6C63FF',
}: RadarChartProps) {
  if (!data || data.length < 3) {
    return (
      <View style={[styles.container, { width: size, height: size }]}>
        <Text style={styles.placeholder}>Need 3+ skills for chart</Text>
      </View>
    );
  }

  const cx = size / 2;
  const cy = size / 2;
  const radius = size * 0.35; // Leave room for labels
  const angleStep = (2 * Math.PI) / data.length;
  const startAngle = -Math.PI / 2; // Start from top

  // Get point on the chart given angle and distance from center
  const getPoint = (angle: number, distance: number) => ({
    x: cx + distance * Math.cos(angle),
    y: cy + distance * Math.sin(angle),
  });

  // Build grid rings
  const gridRings = [];
  for (let level = 1; level <= LEVELS; level++) {
    const r = (radius / LEVELS) * level;
    const points = data
      .map((_, i) => {
        const angle = startAngle + i * angleStep;
        const p = getPoint(angle, r);
        return `${p.x},${p.y}`;
      })
      .join(' ');
    gridRings.push(
      <Polygon
        key={`ring-${level}`}
        points={points}
        fill="none"
        stroke="rgba(148, 163, 184, 0.2)"
        strokeWidth={1}
      />
    );
  }

  // Build axis lines
  const axisLines = data.map((_, i) => {
    const angle = startAngle + i * angleStep;
    const end = getPoint(angle, radius);
    return (
      <Line
        key={`axis-${i}`}
        x1={cx}
        y1={cy}
        x2={end.x}
        y2={end.y}
        stroke="rgba(148, 163, 184, 0.15)"
        strokeWidth={1}
      />
    );
  });

  // Build data polygon
  const dataPoints = data
    .map((d, i) => {
      const angle = startAngle + i * angleStep;
      const distance = (d.value / 100) * radius;
      const p = getPoint(angle, distance);
      return `${p.x},${p.y}`;
    })
    .join(' ');

  // Build data point dots
  const dataDots = data.map((d, i) => {
    const angle = startAngle + i * angleStep;
    const distance = (d.value / 100) * radius;
    const p = getPoint(angle, distance);
    return (
      <Circle
        key={`dot-${i}`}
        cx={p.x}
        cy={p.y}
        r={3}
        fill={RadarAxisColors[i % RadarAxisColors.length]}
        stroke="#fff"
        strokeWidth={1}
      />
    );
  });

  // Build labels
  const labels = data.map((d, i) => {
    const angle = startAngle + i * angleStep;
    const labelDistance = radius + 18;
    const p = getPoint(angle, labelDistance);

    // Determine text anchor based on position
    let textAnchor: 'start' | 'middle' | 'end' = 'middle';
    if (p.x < cx - 5) textAnchor = 'end';
    else if (p.x > cx + 5) textAnchor = 'start';

    // Truncate long labels
    const label = d.label.length > 10 ? d.label.substring(0, 9) + '…' : d.label;

    return (
      <SvgText
        key={`label-${i}`}
        x={p.x}
        y={p.y + 4}
        textAnchor={textAnchor}
        fontSize={10}
        fontWeight="600"
        fill={RadarAxisColors[i % RadarAxisColors.length]}
      >
        {label}
      </SvgText>
    );
  });

  return (
    <View style={[styles.container, { width: size, height: size }]}>
      <Svg width={size} height={size}>
        {gridRings}
        {axisLines}
        <Polygon
          points={dataPoints}
          fill={fillColor}
          stroke={strokeColor}
          strokeWidth={2}
          strokeLinejoin="round"
        />
        {dataDots}
        {labels}
      </Svg>
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    alignItems: 'center',
    justifyContent: 'center',
  },
  placeholder: {
    color: '#94A3B8',
    fontSize: 12,
    textAlign: 'center',
  },
});
