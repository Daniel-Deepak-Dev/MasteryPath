import React from 'react';
import { View, StyleSheet } from 'react-native';
import { Svg, Rect, Text as SvgText } from 'react-native-svg';
import { Colors } from '@/constants/theme';
import { useColorScheme } from '@/hooks/useColorScheme';

interface ContributionDay {
  date: string;
  count: number;
  level: number;
}

interface ContributionGraphProps {
  data: ContributionDay[];
  endDate?: Date;
  days?: number;
}

export const ContributionGraph: React.FC<ContributionGraphProps> = ({
  data,
  endDate = new Date(),
  days = 364, // 52 weeks * 7 days
}) => {
  const colorScheme = useColorScheme();
  const theme = Colors[colorScheme ?? 'light'];
  
  // Configuration
  const squareSize = 12;
  const squareGap = 3;
  const weekLabelWidth = 0; // Hide week labels for now to save space
  const monthLabelHeight = 15;
  const totalWeeks = Math.ceil(days / 7);

  // Helper to get color based on level
  const getColor = (level: number) => {
    switch (level) {
      case 0: return colorScheme === 'dark' ? '#161b22' : '#ebedf0'; // GitHub-like empty
      case 1: return '#9be9a8'; // Light green
      case 2: return '#40c463';
      case 3: return '#30a14e';
      case 4: return '#216e39'; // Dark green
      default: return colorScheme === 'dark' ? '#161b22' : '#ebedf0';
    }
  };

  // Helper to find data for a specific date
  const getDataForDate = (dateStr: string) => {
    return data.find(d => d.date === dateStr) || { date: dateStr, count: 0, level: 0 };
  };

  // Generate grid data
  const weeks = [];
  const monthLabels = [];
  
  let currentMonth = -1;

  for (let w = 0; w < totalWeeks; w++) {
    const weekData = [];
    for (let d = 0; d < 7; d++) {
      const dayOffset = (totalWeeks - 1 - w) * 7 + (6 - d);
      const date = new Date(endDate);
      date.setDate(date.getDate() - dayOffset);
      
      const dateStr = date.toISOString().split('T')[0];
      const dayData = getDataForDate(dateStr);
      weekData.push(dayData);

      // Month labels logic (simplified)
      if (d === 0 && w < totalWeeks - 1) {
         if (date.getMonth() !== currentMonth) {
            currentMonth = date.getMonth();
            monthLabels.push({
               x: w * (squareSize + squareGap),
               label: date.toLocaleString('default', { month: 'short' })
            });
         }
      }
    }
    weeks.push(weekData.reverse()); // Sun -> Sat
  }

  // Calculate dimensions
  const width = totalWeeks * (squareSize + squareGap) + weekLabelWidth;
  const height = 7 * (squareSize + squareGap) + monthLabelHeight;

  return (
    <View style={styles.container}>
      <Svg height={height} width="100%" viewBox={`0 0 ${width} ${height}`}>
        {/* Month Labels */}
        {monthLabels.map((m, i) => (
          <SvgText
            key={`month-${i}`}
            x={m.x}
            y={10}
            fontSize="10"
            fill={theme.text}
            opacity={0.7}
          >
            {m.label}
          </SvgText>
        ))}

        {/* The Grid */}
        {weeks.map((week, wIndex) => (
          week.map((day, dIndex) => (
            <Rect
              key={day.date}
              x={wIndex * (squareSize + squareGap)}
              y={monthLabelHeight + dIndex * (squareSize + squareGap)}
              width={squareSize}
              height={squareSize}
              fill={getColor(day.level)}
              rx={2}
              ry={2}
            />
          ))
        ))}
      </Svg>
    </View>
  );
};

const styles = StyleSheet.create({
  container: {
    paddingVertical: 10,
    alignItems: 'center',
  },
});
